package tg

// ChecklistTask describes a task in a checklist.
type ChecklistTask struct {
	ID              ID              `json:"id"`
	Text            string          `json:"text"`
	Entities        []MessageEntity `json:"text_entities"`
	CompletedByUser *User           `json:"completed_by_user"`
	CompletionDate  Date            `json:"completion_date"`
}

// Checklist describes a checklist.
type Checklist struct {
	Title                    string          `json:"title"`
	Entities                 []MessageEntity `json:"title_entities"`
	Tasks                    []ChecklistTask `json:"tasks"`
	OthersCanAddTasks        bool            `json:"others_can_add_tasks"`
	OthersCanMarkTasksAsDone bool            `json:"others_can_mark_tasks_as_done"`
}

// InputChecklistTask describes a task to add to a checklist.
type InputChecklistTask struct {
	ID        ID              `json:"id"`
	Text      string          `json:"text"`
	ParseMode ParseMode       `json:"parse_mode"`
	Entities  []MessageEntity `json:"text_entities"`
}

// InputChecklist describes a checklist to create.
type InputChecklist struct {
	Title                    string               `json:"title"`
	ParseMode                ParseMode            `json:"parse_mode"`
	Entities                 []MessageEntity      `json:"title_entities"`
	Tasks                    []InputChecklistTask `json:"tasks"`
	OthersCanAddTasks        bool                 `json:"others_can_add_tasks"`
	OthersCanMarkTasksAsDone bool                 `json:"others_can_mark_tasks_as_done"`
}

// ChecklistTasksDone describes a service message about checklist tasks marked as done or not done.
type ChecklistTasksDone struct {
	Message                *Message `json:"checklist_message"`
	MarkedAsDoneTaskIDs    []ID     `json:"marked_as_done_task_ids"`
	MarkedAsNotDoneTaskIDs []ID     `json:"marked_as_not_done_task_ids"`
}

// ChecklistTasksAdded describes a service message about tasks added to a checklist.
type ChecklistTasksAdded struct {
	Message *Message        `json:"checklist_message"`
	Tasks   []ChecklistTask `json:"tasks"`
}
