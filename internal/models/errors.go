package models

import "errors"

var (
	ErrCourseNotFound     = errors.New("course not found")
	ErrLessonNotFound     = errors.New("lesson not found")
	ErrUserNotFound       = errors.New("user not found")
	ErrTeacherNotFound    = errors.New("teacher not found")
	ErrSlugAlreadyExists  = errors.New("course slug already exists")
	ErrEmailAlreadyExists = errors.New("user email already exists")
)
