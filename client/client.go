package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	pb "student-analytics/proto"
	"time"

	"google.golang.org/grpc"
)

// main is the entry point of the client application. It connects to the gRPC server,
// interacts with the student service, and handles user input for adding new students,
// retrieving student details, and calculating the average grade.
func main() {
	// Connect to the gRPC server
	conn, err := grpc.Dial("localhost:3000", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewStudentServiceClient(conn)

	// Increase timeout for the context
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second) // Increased timeout
	defer cancel()

	// Initial input for a new student
	var name string
	var age int
	var grade float64
	var gender string

	// Prompt user for input
	fmt.Println("Enter student details:")
	fmt.Print("Name: ")
	fmt.Scanln(&name)

	fmt.Print("Age: ")
	ageStr := ""
	fmt.Scanln(&ageStr)
	age, err = strconv.Atoi(ageStr)
	if err != nil {
		log.Fatalf("Invalid age format: %v", err)
	}

	fmt.Print("Grade: ")
	gradeStr := ""
	fmt.Scanln(&gradeStr)
	grade, err = strconv.ParseFloat(gradeStr, 32)
	if err != nil {
		log.Fatalf("Invalid grade format: %v", err)
	}

	fmt.Print("Gender ( M | F  )  ? : ")
	fmt.Scanln(&gender)

	// Add the student to the database
	_, err = client.AddStudent(ctx, &pb.Student{Name: name, Age: int32(age), Grade: float32(grade), Gender: gender})
	if err != nil {
		log.Fatalf("could not add student: %v", err)
	}
	fmt.Println("Student added successfully.")

	// Retrieve a student by name
	fmt.Print("Enter the name of the student to retrieve: ")
	fmt.Scanln(&name)

	student, err := client.GetStudent(ctx, &pb.StudentRequest{Name: name})
	if err != nil {
		log.Fatalf("could not get student: %v", err)
	}

	log.Printf("Retrieved Student: Name: %s, Age: %d, Grade: %.2f", student.Name, student.Age, student.Grade)

	// Get average grade
	avgResp, err := client.GetAverageGrade(ctx, &pb.Empty{})
	if err != nil {
		log.Fatalf("could not get average grade: %v", err)
	}
	log.Printf("Average Grade: %f", avgResp.AverageGrade)

	// getmax age by gender
	maxResp, err := client.GetMaxAgeByGender(ctx, &pb.Empty{})
	if err != nil {
		log.Fatalf("could not get max age by gender: %v", err)
	}
	log.Printf("Max Age by Gender: Male: %d, Female: %d", maxResp.MaleMaxAge, maxResp.FemaleMaxAge)

	// getmin age by gender

	minResp, err := client.GetMinAgeByGender(ctx, &pb.Empty{})
	if err != nil {
		log.Fatalf("could not get min age by gender: %v", err)
	}
	log.Printf("Min Age by Gender: Male: %d, Female: %d", minResp.MaleMinAge, minResp.FemaleMinAge)

	// get gender precenatge

	genderResp, err := client.GetGenderPercentage(ctx, &pb.Empty{})
	if err != nil {
		log.Fatalf("could not get gender percentage: %v", err)
	}
	log.Printf("Gender Percentage: Male: %.2f, Female: %.2f", 100*genderResp.MalePercentage, 100*genderResp.FemalePercentage)

	//get averagfe grade by gender

	avgGradeResp, err := client.GetAverageGradeByGender(ctx, &pb.Empty{})
	if err != nil {
		log.Fatalf("could not get average grade by gender: %v", err)
	}
	log.Printf("Average Grade by Gender: Male: %.2f, Female: %.2f", avgGradeResp.MaleAverageGrade, avgGradeResp.FemaleAverageGrade)

	empty := &pb.Empty{}

	// Call the GetCombinedData method on the server
	combinedResponse, err := client.GetCombinedData(context.Background(), empty)
	if err != nil {
		log.Fatalf("failed to get combined data: %v", err)
	}

	// Marshal the response data into JSON
	jsonBytes, err := json.Marshal(combinedResponse)
	if err != nil {
		log.Fatalf("failed to marshal combined data: %v", err)
	}

	// Print the JSON data to the terminal
	fmt.Println(string(jsonBytes))
}
