package notifier

import (
	"github.com/reactivejson/users-svc/internal/domain"
)

/**
 * @author Mohamed-Aly Bou-Hanane
 * © 2023
 */

// UserNotifier represents the user notifier.
type UserNotifier interface {
	NotifyUserChange(eventType UserEventType, user *domain.User) error
}
