package winutil

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows"
)

// AdminOnlySDDL grants full control to SYSTEM and Administrators only,
// with inheritance from the parent blocked (protected DACL).
const AdminOnlySDDL = "D:P(A;OICI;FA;;;SY)(A;OICI;FA;;;BA)"

// SecureDir creates dir if needed and replaces its DACL with
// AdminOnlySDDL, so programs running as a normal user can neither read
// nor replace what VinPN (elevated) keeps there. It also takes ownership
// for Administrators and removes children that are not admin-owned: a
// process may have guessed the path and planted files under it while the
// parent (%ProgramData%) still allowed folder creation, and an owner keeps
// write rights over any DACL. Children VinPN wrote itself are always
// admin-owned (the app runs elevated), so only planted entries disappear.
func SecureDir(dir string) error {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	owned, _ := OwnedByAdmins(dir)
	harden := owned
	if !owned {
		if err := takeOwnership(dir); err == nil {
			harden = true
		}
		// Not owner-transferable (e.g. a non-elevated dev run on a
		// directory it created itself): the DACL below still applies;
		// sweeping would delete the caller's own files, so skip it.
	}
	sd, err := windows.SecurityDescriptorFromString(AdminOnlySDDL)
	if err != nil {
		return err
	}
	dacl, _, err := sd.DACL()
	if err != nil {
		return err
	}
	err = windows.SetNamedSecurityInfo(dir, windows.SE_FILE_OBJECT,
		windows.DACL_SECURITY_INFORMATION|windows.PROTECTED_DACL_SECURITY_INFORMATION, nil, nil, dacl, nil)
	if err != nil {
		return fmt.Errorf("secure %s: %w", dir, err)
	}
	if harden {
		return sweepUnowned(dir)
	}
	return nil
}

// takeOwnership hands dir to the built-in Administrators group. Elevated
// VinPN holds SeTakeOwnership, which is what makes this work against a
// directory a normal user created first.
func takeOwnership(dir string) error {
	ba, err := windows.CreateWellKnownSid(windows.WinBuiltinAdministratorsSid)
	if err != nil {
		return err
	}
	if err := windows.SetNamedSecurityInfo(dir, windows.SE_FILE_OBJECT,
		windows.OWNER_SECURITY_INFORMATION, ba, nil, nil, nil); err != nil {
		return fmt.Errorf("secure %s: owner: %w", dir, err)
	}
	return nil
}

// sweepUnowned removes direct children not owned by Administrators or
// SYSTEM. The just-applied DACL (FA for SY and BA, DELETE_CHILD included)
// lets the elevated caller remove them regardless of the child's own DACL.
func sweepUnowned(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	var errs []error
	for _, e := range entries {
		p := filepath.Join(dir, e.Name())
		if ok, err := OwnedByAdmins(p); err == nil && !ok {
			if err := os.RemoveAll(p); err != nil {
				errs = append(errs, fmt.Errorf("secure %s: remove %s: %w", dir, e.Name(), err))
			}
		}
	}
	return errors.Join(errs...)
}

// OwnedByAdmins reports whether path is owned by Administrators or SYSTEM.
// Files VinPN writes while elevated are; a file planted by a normal
// user is not.
func OwnedByAdmins(path string) (bool, error) {
	sd, err := windows.GetNamedSecurityInfo(path, windows.SE_FILE_OBJECT, windows.OWNER_SECURITY_INFORMATION)
	if err != nil {
		return false, err
	}
	owner, _, err := sd.Owner()
	if err != nil {
		return false, err
	}
	for _, w := range []windows.WELL_KNOWN_SID_TYPE{windows.WinBuiltinAdministratorsSid, windows.WinLocalSystemSid} {
		sid, err := windows.CreateWellKnownSid(w)
		if err == nil && owner.Equals(sid) {
			return true, nil
		}
	}
	return false, nil
}
