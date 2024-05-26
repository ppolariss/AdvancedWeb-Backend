package exam

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(routes fiber.Router) {
	// exam
	routes.Get("/exams", ListExams)
	routes.Get("/exams/:id", GetExam)
	routes.Post("/exams/start", StartExam)
	routes.Post("/exams/end", EndExam)

	routes.Post("/exams/add", AddExam)
	routes.Put("/exams", ModifyExam)
	routes.Delete("/exams/:id", DeleteExam)

	// exam punishment
	routes.Get("/exams/punishments/:id", GetPunishment)
	routes.Post("/exams/:id/punishments/", AddPunishment)
	routes.Get("/exams/:id<int>/punishments/", ListPunishments)
	//routes.Get("/exams/:id/questions", ListQuestions)
	//routes.Post("/exams/:id/questions", AddQuestion)
	//routes.Put("/questions/:id", ModifyQuestion)
	//routes.Delete("/questions/:id", DeleteQuestion)
	//routes.Get("/questions/:id/choices", ListChoices)
	//routes.Post("/questions/:id/choices", AddChoice)
	//routes.Put("/choices/:id", ModifyChoice)
	//routes.Delete("/choices/:id", DeleteChoice)
	//routes.Get("/questions/:id/answers", ListAnswers)
	//routes.Post("/questions/:id/answers", AddAnswer)
	//routes.Put("/answers/:id", ModifyAnswer)
	//routes.Delete("/answers/:id", DeleteAnswer)

	routes.Get("/society/:id", GetSocietyPunishment)
	routes.Post("/society/punishments/", AddSocietyPunishment)
	routes.Get("/society/punishments/", ListSocietyPunishments)
}
