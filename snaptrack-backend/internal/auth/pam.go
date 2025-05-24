package auth

import (
    "errors"
    "fmt"
    "os/exec"
    "strings"

    "github.com/msteinert/pam"
)

// AuthenticateUser authenticates username/password using PAM and checks sudo/root group
func AuthenticateUser(username, password string) error {
    tx, err := pam.StartFunc("login", username, func(style pam.Style, msg string) (string, error) {
        return password, nil
    })
    if err != nil {
        return err
    }
    if err := tx.Authenticate(0); err != nil {
        return err
    }

    // If root, allow
    if username == "root" {
        return nil
    }

    // Check if in sudo or wheel group
    out, err := exec.Command("groups", username).Output()
    if err != nil {
        return fmt.Errorf("failed to get groups for user %s: %w", username, err)
    }

    groups := strings.Fields(string(out))
    if len(groups) < 3 {
        return errors.New("unable to parse groups for user")
    }
    userGroups := groups[2:]
    for _, g := range userGroups {
        if g == "sudo" || g == "wheel" {
            return nil
        }
    }

    return errors.New("user does not have root or sudo privileges")
}
