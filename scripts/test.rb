#!/usr/bin/env ruby
# coding: utf-8

# Run unit and integration tests for preservation-services.

require 'fileutils'
require 'net/http'
require 'optparse'

class TestRunner

  attr_accessor :test_name

  def initialize(options)
    @options = options
    @pids = {}
    @services_stopped = false
    @test_name = '';
    @start_time = Time.now
    bin = self.bin_dir
    @integration_services = [
      {
        # For localhost testing, use 'localhost' instead of '127.0.0.1'
        # because Minio signed URLs use hostname, not IP.
        name: "minio",
        cmd: "#{bin}/minio server --quiet --address=localhost:9899 ~/tmp/minio",
        msg: "Minio is running on localhost:9899. User/Pwd: minioadmin/minioadmin"
      }
    ]
  end

  def clean_test_cache
    puts "Deleting test cache from last run"
    `go clean -testcache`
    puts "Deleting old Redis data"
    File.delete('dump.rdb') if File.exists?('dump.rdb')
  end


  def run_unit_tests(arg)
    clean_test_cache
    run_go_unit_tests(arg)
  end

  def run_go_unit_tests(arg)
    puts "Starting unit tests..."
    arg = "./..." if arg.nil?
    cmd = "go test #{arg}"
    puts cmd
    pid = Process.spawn(env_hash, cmd, chdir: project_root)
    Process.wait pid
    self.print_results
  end

  def run_integration_tests(arg)
    init_for_integration
    puts "Starting integration tests..."
    arg = "./..." if arg.nil?
    cmd = "go test -tags=integration #{arg}"
    puts cmd
    pid = Process.spawn(env_hash, cmd, chdir: project_root)
    Process.wait pid
    self.print_results
  end


  # Initialize for integration, interactive tests, and
  # end to end tests. This clears and rebuilds data directories,
  # starts all services, and creates all NSQ topics.
  def init_for_integration
    clean_test_cache
    make_test_dirs
    self.registry_start
    sleep(8)
    @integration_services.each do |svc|
      start_service(svc)
    end
    sleep(5)
  end

  def start_service(svc)
    log_file = log_file_path(svc[:name])
    pid = Process.spawn(env_hash, svc[:cmd], out: log_file, err: log_file)
    Process.detach pid
    log_started(svc, pid, log_file)
  	@pids[svc[:name]] = pid
  end

  def log_started(svc, pid, log_file)
    puts ""
    puts "Started #{svc[:name]} with command '#{svc[:cmd]}' and pid #{pid}"
    puts svc[:msg]
    puts "Log file is #{log_file}"
    puts ""
  end

  def stop_service(name, pid)
	if pid.nil? || pid == 0
      puts "Pid for #{name} is zero. Can't kill that..."
	  return
	end
    os = (/darwin/ =~ RUBY_PLATFORM) ? "osx" : "linux"
    if os == "linux"
      stop_service_linux(name)
      return
    end
    puts "Stopping #{name} service (pid #{pid})"
    begin
  	  Process.kill('TERM', pid)
  	rescue
	  puts "Hmm... Couldn't kill #{name}."
      puts "Check system processes to see if a version "
      puts "of that process is lingering from a previous test run."
	end
  end

  # This method exists because Process.spawn on Linux returns the
  # pid of a short-lived parent process, which creates the service
  # and then exits. That means we can't know the pid of the actual
  # service we want to kill.
  #
  # Note that killing a process by name carries some risk. This will
  # kill ALL nsq, redis, minio, and registry processes. That should
  # be OK on dev/test systems, but if you're wondering where your
  # redis/minio/nsq/registry process went, it went down this drain.
  #
  # If you're running these tests in a system that has its own long-
  # running minio/redis/nsq services, the tests will likely fail
  # anyway because those services will hold on to stale data.
  def stop_service_linux(name)
    pids = `pidof #{name}`
    pids.split(' ').each do |pid|
      begin
        Process.kill('TERM', pid.to_i)
        puts "(Linux) Killed #{name} with pid #{pid}"
      rescue
        puts "Hmm... Couldn't kill #{name}."
        puts "Check system processes to see if a version "
        puts "of that process is lingering from a previous test run."
      end
    end
  end

  def env_hash
    env = {}
	  ENV.each{ |k,v| env[k] = v }
	  # env['APT_ENV'] = 'integration'
    if self.test_name != 'units'
      env['REGISTRY_ROOT'] = ENV['REGISTRY_ROOT'] || abort("Set env var REGISTRY_ROOT")
    end
	  env
  end

  def make_test_dirs
    base = File.join(ENV['HOME'], "tmp")
    if base.end_with?("tmp") # So we don't delete anyone's home dir
      puts "Deleting #{base}"
    end
    FileUtils.remove_dir(base ,true)
    dirs = ["bin", "logs", "minio", "nsq", "redis", "restore"]
    dirs.each do |dir|
      full_dir = File.join(base, dir)
      puts "Creating #{full_dir}"
      FileUtils.mkdir_p full_dir
    end
    # S3 buckets for minio. We should ideally read these from the
    # .env.test file.
    buckets = [
      "test-bucket-1",
      "test-bucket-2",
    ]
    buckets.each do |bucket|
      full_bucket = File.join(base, "minio", bucket)
      puts "Creating local minio bucket #{bucket}"
      FileUtils.mkdir_p full_bucket
    end
  end

  def project_root
    File.expand_path(File.join(File.dirname(__FILE__), ".."))
  end

  def ingest_bin_dir
    File.join(project_root, "bin", "go-bin")
  end

  def bin_dir
    os = (/darwin/ =~ RUBY_PLATFORM) ? "osx" : "linux"
    File.join(project_root, "bin", os)
  end

  # Note: This assumes you have the registry repo source tree
  # on your machine. It's on GitHub at https://github.com/APTrust/registry
  def registry_start
  	if !@pids['registry']
      registry_load_fixtures
	    # Force copy of env to integration so that registry fixtures load.
	    env = {}.merge(env_hash)
	    env['APT_ENV'] = 'integration'
      # Important! Adding -tags=test here turns on the special
      # testing endpoints prepare_file_delete and prepare_object_delete,
      # which are disabled in all non-test environments.
	    cmd = 'go run -tags=test registry.go'
	    log_file = log_file_path('registry')
	    registry_pid = Process.spawn(env,
								   cmd,
								   chdir: env['REGISTRY_ROOT'],
								   out: [log_file, 'w'],
								   err: [log_file, 'w'])
	    Process.detach registry_pid
      sleep 3

      # go run compiles an executable, puts it in a temp directory, and
      # runs it as a new process. We need to get the pid of that process.
      # Note that the temp dir pattern will be different on linux.
      # /var/folders works for Mac.
      registry_process = `ps -ef | grep registry | grep /var/folders`
      pid = registry_process.split(/\s+/)[2].to_i
      if pid
        @pids['registry'] = pid
      else
        @pids['registry'] = registry_pid
      end
	    puts "Started Registry with command '#{cmd}' and pid #{@pids['registry']}"
	  end
  end

  def registry_load_fixtures
  	puts "Loading registry fixtures"
	  env = {}.merge(env_hash)
	  env['APT_ENV'] = 'integration'
	  cmd = 'go run loader/load_fixtures.go'
	  log_file = log_file_path('registry_fixtures')
	  registry_pid = Process.spawn(env,
								 cmd,
								 chdir: env['REGISTRY_ROOT'],
								 out: [log_file, 'w'],
								 err: [log_file, 'w'])
	  Process.wait
    puts "Registry fixtures loaded"
  end

  def log_file_path(service_name)
    return File.join(ENV['HOME'], "tmp", "logs", service_name + ".log")
  end

  def stop_all_services
    return if @services_stopped
    puts "Stopping all services"
    @pids.each do |name, pid|
      stop_service(name, pid)
    end
    @services_stopped = true
  end

  def print_results
    puts "\n"
    puts "Elapsed time: #{Time.now - @start_time} seconds"
    puts "Logs are in #{File.join(ENV['HOME'], "tmp", "logs")}"
    if $?.success?
      puts "\n\n    **** 😁 PASS 😁 **** \n\n".force_encoding('utf-8')
    else
      puts "\n\n    **** 🤬 FAIL 🤬 **** \n\n".force_encoding('utf-8')
      exit(false)
    end
  end

  def print_help
    puts "\n"
    puts "APTrust partner tools tests\n\n"
	  puts "Usage: "
    puts "  test.rb units                   # Run unit tests"
    puts "  test.rb integration             # Run integration tests"
    puts "\n"
    puts "Note that running integration tests also runs unit tests."
    puts "\n"
  end

end

# TODO: Add command line args to specify whether to run unit tests
# or integration tests. For now, we're only running unit tests.
if __FILE__ == $0
  options = {}

  t = TestRunner.new(options)
  t.test_name = ARGV[0]
  if !['units', 'integration'].include?(t.test_name)
    t.print_help
	exit(false)
  end
  at_exit { t.stop_all_services }
  case t.test_name
  when 'units'
    t.run_unit_tests(ARGV[1])
  when 'integration'
    t.run_integration_tests(ARGV[1])
  end
end
