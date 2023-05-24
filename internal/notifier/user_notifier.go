// internal/notifier/user_notifier.go
package notifier

import (
	"github.com/reactivejson/usr-svc/internal/domain"
)

// UserNotifier represents the user notifier.
type UserNotifier interface {
	NotifyUserChange(user *domain.User) error
}
