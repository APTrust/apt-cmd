package cmd

const (
	// EXIT_OK means program completed successfully.
	EXIT_OK = 0
	// EXIT_RUNTIME_ERR means program did not complete
	// successfully due to an error. The error may have
	// occurred outside the program, such as a network
	// error or an error on a remote server.
	EXIT_RUNTIME_ERR = 1
	// EXIT_BAG_INVALID is used primarily for apt_validate.
	// It means the program completed its run and found that
	// the bag is not valid.
	EXIT_BAG_INVALID = 2
	// EXIT_USER_ERR means the user did not supply some
	// required option or argument, or the user supplied
	// invalid options/arguments.
	EXIT_USER_ERR = 3
	// EXIT_NO_OP means the user requested help message or
	// version info. The program printed the info, and no other
	// operations were performed.
	EXIT_NO_OP = 100
)
