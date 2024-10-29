package tasks

import (
	orderTasks "haraka-sana/orders/tasks"
)

func ListenEvents() {
	orderTasks.RelayOrderEvents()
}
