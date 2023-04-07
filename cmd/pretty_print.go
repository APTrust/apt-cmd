package cmd

import (
	"fmt"
	"time"

	"github.com/APTrust/preservation-services/constants"
	"github.com/APTrust/preservation-services/models/registry"
	"github.com/APTrust/preservation-services/network"
)

func PrettyPrint(resp *network.RegistryResponse, objType network.RegistryObjectType) {
	fmt.Println(objType, ":", resp.Count, "results")
	switch objType {
	case network.RegistryIntellectualObject:
		prettyPrintObjectList(resp.IntellectualObjects())
	case network.RegistryGenericFile:
		prettyPrintFileList(resp.GenericFiles())
	case network.RegistryWorkItem:
		prettyPrintWorkItemList(resp.WorkItems())
	}
	if resp.HasNextPage() {
		fmt.Println("Next Page:", resp.Next)
	}
}

func prettyPrintObjectList(list []*registry.IntellectualObject) {
	for _, obj := range list {
		fmt.Println("")
		PrettyPrintObject(obj)
		fmt.Println("======================================================")
	}
}

func prettyPrintFileList(list []*registry.GenericFile) {
	for _, gf := range list {
		fmt.Println("")
		PrettyPrintFile(gf)
		fmt.Println("======================================================")
	}
}

func prettyPrintWorkItemList(list []*registry.WorkItem) {
	for _, item := range list {
		fmt.Println("")
		PrettyPrintWorkItem(item)
		fmt.Println("======================================================")
	}
}

func PrettyPrintObject(obj *registry.IntellectualObject) {

	if obj == nil || obj.ID == 0 {
		fmt.Println("No results")
		return
	}

	state := "Active"
	if obj.State == constants.StateDeleted {
		state = "Deleted"
	}
	bagItProfile := "APTrust"
	if obj.BagItProfileIdentifier == constants.BTRProfileIdentifier {
		bagItProfile = "BTR"
	}

	// Identifier Info
	fmt.Println("ID:            ", obj.ID)
	fmt.Println("Identifier:    ", obj.Identifier)
	fmt.Println("Bag Name:      ", obj.BagName)
	fmt.Println("ETag:          ", obj.ETag)
	fmt.Println("Inst ID:       ", obj.InstitutionID)
	fmt.Println("Institution:   ", obj.InstitutionIdentifier)
	fmt.Println("Title:         ", obj.Title)
	fmt.Println("Access:        ", obj.Access)
	fmt.Println("State:         ", state)
	fmt.Println("Storage Opt:   ", obj.StorageOption)
	fmt.Println("BagIt Profile: ", bagItProfile)

	// Stats
	fmt.Println("--- Object Stats ---")
	fmt.Println("Size:          ", obj.Size)
	fmt.Println("File Count:    ", obj.FileCount)
	fmt.Println("Payload Files: ", obj.PayloadFileCount)
	fmt.Println("Payload Size:  ", obj.PayloadSize)
	fmt.Println("Created At:    ", obj.CreatedAt)
	fmt.Println("Updated At:    ", obj.UpdatedAt)

	// Inst internal info
	fmt.Println("--- Internal Info ---")
	fmt.Println("Alt Identifier:       ", obj.AltIdentifier)
	fmt.Println("Bag Group:            ", obj.BagGroupIdentifier)
	fmt.Println("Source Organization:  ", obj.SourceOrganization)
	fmt.Println("Internal Identifier:  ", obj.InternalSenderIdentifier)
	fmt.Println("Internal Description: ", obj.InternalSenderDescription)

	// Show description last, as it can be lengthy
	fmt.Println("External Description: ", obj.Description)
	fmt.Println("")
}

func PrettyPrintFile(gf *registry.GenericFile) {

	if gf == nil || gf.ID == 0 {
		fmt.Println("No results")
		return
	}

	fmt.Println("File ID:     ", gf.ID)
	fmt.Println("Identifier:  ", gf.Identifier)
	fmt.Println("UUID:        ", gf.UUID)
	fmt.Println("Size:        ", gf.Size)
	fmt.Println("State:       ", gf.State)
	fmt.Println("Format:      ", gf.FileFormat)
	fmt.Println("Storage Opt: ", gf.StorageOption)
	fmt.Println("Inst ID:     ", gf.InstitutionID)
	fmt.Println("Object ID:   ", gf.IntellectualObjectID)
	fmt.Println("Created At:  ", gf.CreatedAt)
	fmt.Println("Updated At:  ", gf.UpdatedAt)
	fmt.Println("Modified:    ", gf.FileModified)
	fmt.Println("Last Fixity: ", gf.LastFixityCheck)

	fmt.Println("")
	fmt.Println("--- Storage Records ---")
	for _, sr := range gf.StorageRecords {
		PrettyPrintStorageRecord(sr)
	}

	fmt.Println("")
	fmt.Println("--- Premis Events ---")
	for _, event := range gf.PremisEvents {
		PrettyPrintEvent(event)
	}

	fmt.Println("")
	fmt.Println("--- Checksums ---")
	for _, cs := range gf.Checksums {
		PrettyPrintChecksum(cs)
	}
	fmt.Println("")
}

func PrettyPrintChecksum(cs *registry.Checksum) {
	if cs == nil || cs.ID == 0 {
		fmt.Println("No results")
		return
	}
	fmt.Println("Date:      ", cs.DateTime.Format(time.RFC3339))
	fmt.Println("Algorithm: ", cs.Algorithm)
	fmt.Println("Digest:    ", cs.Digest)
	fmt.Println("")
}

func PrettyPrintEvent(event *registry.PremisEvent) {
	if event == nil || event.ID == 0 {
		fmt.Println("No results")
		return
	}

	fmt.Println("Event ID:       ", event.ID)
	fmt.Println("UUID:           ", event.Identifier)
	fmt.Println("Date:           ", event.DateTime)
	fmt.Println("Type:           ", event.EventType)
	fmt.Println("Outcome:        ", event.Outcome)
	fmt.Println("Agent:          ", event.Agent)
	fmt.Println("Detail:         ", event.Detail)
	fmt.Println("Object:         ", event.Object)
	fmt.Println("Outcome Detail: ", event.OutcomeDetail)
	fmt.Println("Info:           ", event.OutcomeInformation)
	fmt.Println("")
}

func PrettyPrintStorageRecord(sr *registry.StorageRecord) {
	if sr == nil || sr.ID == 0 {
		fmt.Println("No results")
		return
	}
	fmt.Println("ID:  ", sr.ID)
	fmt.Println("URL: ", sr.URL)
	fmt.Println("")
}

func PrettyPrintWorkItem(item *registry.WorkItem) {
	if item == nil || item.ID == 0 {
		fmt.Println("No results")
		return
	}

	// General: common to all items
	fmt.Println("Work Item ID:   ", item.ID)
	fmt.Println("Name:           ", item.Name)
	fmt.Println("Size:           ", item.Size)
	fmt.Println("Action:         ", item.Action)
	fmt.Println("Stage:          ", item.Stage)
	fmt.Println("Status:         ", item.Status)
	fmt.Println("Needs Review:   ", item.NeedsAdminReview)
	fmt.Println("Updated At:     ", item.UpdatedAt)
	fmt.Println("Note:           ", item.Note)

	// Obj Info
	fmt.Println("--- Object Info ---")
	fmt.Println("Object ID:      ", item.IntellectualObjectID)
	fmt.Println("Obj Identifier: ", item.ObjectIdentifier)
	fmt.Println("Bag Date:       ", item.BagDate.Format(time.RFC3339))
	fmt.Println("ETag:           ", item.ETag)
	fmt.Println("Storage Option: ", item.StorageOption)

	// File Info
	if item.GenericFileID > 0 {
		fmt.Println("--- File Info ---")
		fmt.Println("File ID:         ", item.GenericFileID)
		fmt.Println("File Identifier: ", item.GenericFileIdentifier)
	}

	// Restorations & deletions
	if item.Action != constants.ActionIngest {
		if item.Action == constants.ActionDelete {
			fmt.Println("--- Deletion Info ---")
		} else {
			// Obj restore, file restore, Glacier restore
			fmt.Println("--- Restoration Info ---")
		}
		fmt.Println("Requested By:    ", item.User)
		fmt.Println("Approved By:     ", item.InstApprover)
	}
	fmt.Println("")
}
