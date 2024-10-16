package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	models "student-analytics/model"
	pb "student-analytics/proto"

	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type server struct {
	pb.UnimplementedStudentServiceServer
	db *gorm.DB
}

func (s *server) AddStudent(ctx context.Context, student *pb.Student) (*pb.Empty, error) {
	newStudent := models.Student{
		Name:   student.Name,
		Age:    int(student.Age),
		Grade:  student.Grade,
		Gender: student.Gender,
	}
	if err := s.db.Create(&newStudent).Error; err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}

func (s *server) GetStudent(ctx context.Context, req *pb.StudentRequest) (*pb.Student, error) {
	var student models.Student
	if err := s.db.First(&student, "name = ?", req.Name).Error; err != nil {
		return nil, err
	}
	return &pb.Student{Name: student.Name, Age: int32(student.Age), Grade: student.Grade}, nil
}

func (s *server) GetAverageGrade(ctx context.Context, empty *pb.Empty) (*pb.AverageGradeResponse, error) {
	var students []models.Student
	var total float64
	s.db.Find(&students)
	for _, student := range students {
		total += float64(student.Grade)
	}
	average := total / float64(len(students))
	return &pb.AverageGradeResponse{AverageGrade: float32(average)}, nil
}

func (s *server) GetGenderPercentage(ctx context.Context, empty *pb.Empty) (*pb.PercentageResponse, error) {

	var students []models.Student
	s.db.Find(&students)

	var maleCount int32
	var femaleCount int32

	for _, student := range students {
		if student.Gender == "M" {
			maleCount++
		} else if student.Gender == "F" {
			femaleCount++
		}
	}

	malePercentage := float32(maleCount) / float32(len(students))
	femalePercentage := float32(femaleCount) / float32(len(students))
	fmt.Println("male precentage : ", malePercentage, "female precentage : ", femalePercentage)
	return &pb.PercentageResponse{
		MalePercentage:   malePercentage,
		FemalePercentage: femalePercentage,
	}, nil

}

func (s *server) GetMaxAgeByGender(ctx context.Context, empty *pb.Empty) (*pb.MaxAgeByGenderResponse, error) {

	var students []models.Student
	var maleMaxAge, femaleMaxAge int

	s.db.Model(&students).Where("Gender = ?", "M").Select("MAX(Age) AS Age").Scan(&maleMaxAge)
	s.db.Model(&students).Where("Gender = ?", "F").Select("MAX(Age) AS Age").Scan(&femaleMaxAge)
	return &pb.MaxAgeByGenderResponse{
		MaleMaxAge:   int32(maleMaxAge),
		FemaleMaxAge: int32(femaleMaxAge),
	}, nil
}

func (s *server) GetMinAgeByGender(ctx context.Context, empty *pb.Empty) (*pb.MinAgeByGenderResponse, error) {
	var students []models.Student
	var maleMinAge, femaleMinAge int

	s.db.Model(&students).Where("Gender = ?", "M").Select("Min(Age) AS Age").Scan(&maleMinAge)
	s.db.Model(&students).Where("Gender = ?", "F").Select("Min(Age) AS Age").Scan(&femaleMinAge)
	return &pb.MinAgeByGenderResponse{
		MaleMinAge:   int32(maleMinAge),
		FemaleMinAge: int32(femaleMinAge),
	}, nil
}

func (s *server) GetAverageGradeByGender(ctx context.Context, empty *pb.Empty) (*pb.AverageGradeByGenderResponse, error) {
	var students []models.Student

	var maleAvg, femaleAvg float64
	s.db.Model(&students).Where("Gender = ?", "M").Select("AVG(Grade) AS Grade").Scan(&maleAvg)
	s.db.Model(&students).Where("Gender = ?", "F").Select("AVG(Grade) AS Grade").Scan(&femaleAvg)

	return &pb.AverageGradeByGenderResponse{
		MaleAverageGrade:   float32(maleAvg),
		FemaleAverageGrade: float32(femaleAvg),
	}, nil

}

func connectDB() (*gorm.DB, error) {

	// Retrieve environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Create the Data Source Name (DSN)
	dsn := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}

	if err := models.Migrate(db); err != nil {
		log.Fatal("failed to migrate database: ", err)
	}

	return db, nil
}

func (s *server) GetCombinedData(ctx context.Context, req *pb.Empty) (*pb.CombinedResponse, error) {
	// Reuse the existing GetAverageGrade function
	avgGradeResp, err := s.GetAverageGrade(ctx, req)
	if err != nil {
		return nil, err
	}

	// Reuse the existing GetGenderPercentage function
	genderPercentageResp, err := s.GetGenderPercentage(ctx, req)
	if err != nil {
		return nil, err
	}

	// Reuse the existing GetMaxAgeByGender function
	maxAgeResp, err := s.GetMaxAgeByGender(ctx, req)
	if err != nil {
		return nil, err
	}

	// Reuse the existing GetMinAgeByGender function
	minAgeResp, err := s.GetMinAgeByGender(ctx, req)
	if err != nil {
		return nil, err
	}

	// Fetch all students to include in the combined response
	var students []models.Student
	if err := s.db.Find(&students).Error; err != nil {
		return nil, err
	}

	// Convert students to protobuf format
	var pbStudents []*pb.Student
	for _, student := range students {
		pbStudents = append(pbStudents, &pb.Student{
			Name:   student.Name,
			Age:    int32(student.Age),
			Grade:  student.Grade,
			Gender: student.Gender,
		})
	}

	// Combine all the results into a single response
	combinedResponse := &pb.CombinedResponse{
		Students:         pbStudents,
		AverageGrade:     avgGradeResp.AverageGrade,
		GenderPercentage: genderPercentageResp,
		MaxAge:           maxAgeResp,
		MinAge:           minAgeResp,
	}

	return combinedResponse, nil
}
func main() {
	srv := grpc.NewServer()

	database, err := connectDB()

	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	pb.RegisterStudentServiceServer(srv, &server{db: database})

	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatalf("failed to listen on port 3000: %v", err)
	}

	if err := srv.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
