package helper

import (
	cfg "ahava/pkg/config"
	"ahava/pkg/utils/models"
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	// "github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/golang-jwt/jwt"
	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"

	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/verify/v2"

	"crypto/rand"
	"encoding/base32"
)

type Helper interface {
	GenerateTokenAdmin(admin models.AdminDetailsResponse) (string, string, error)
	AddFileToS3(file *multipart.FileHeader, bucketName string) (string, error)
	TwilioSetup(username string, password string)
	TwilioSendOTP(phone string, serviceID string) (string, error)
	TwilioVerifyOTP(serviceID string, code string, phone string) error
	GenerateTokenClients(user models.UserDetailsResponse) (string, error)
	GenerateRefferalCode() (string, error)
	PasswordHashing(string) (string, error)
	CompareHashAndPassword(a string, b string) error
	Copy(a *models.UserDetailsResponse, b *models.UserSignInResponse) (models.UserDetailsResponse, error)
}

type helper struct {
	cfg cfg.Config
}

func NewHelper(config cfg.Config) Helper {
	return &helper{
		cfg: config,
	}
}

var client *twilio.RestClient

type AuthCustomClaims struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.StandardClaims
}

func (h *helper) GenerateTokenAdmin(admin models.AdminDetailsResponse) (string, string, error) {
	accessTokenClaims := &AuthCustomClaims{
		ID:    admin.ID,
		Email: admin.Email,
		Role:  "admin",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 6).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	refreshTokenClaims := &AuthCustomClaims{
		ID:    admin.ID,
		Email: admin.Email,
		Role:  "admin",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 30).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString([]byte("accesssecret"))
	if err != nil {
		return "", "", err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte("refreshsecret"))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func (h *helper) AddFileToS3(file *multipart.FileHeader, bucketName string) (string, error) {
	// Initialize MinIO client
	minioClient, err := minio.New(h.cfg.MINIO_ENDPOINT, &minio.Options{
		Creds:  credentials.NewStaticV4(h.cfg.MINIO_ACCESS_KEY_ID, h.cfg.MINIO_SECRET_ACCESS_KEY, ""),
		Secure: false,
	})
	if err != nil {
		log.Printf("Error initializing MinIO client: %v", err)
		return "", err
	}

	// Open file
	fileContent, err := file.Open()
	if err != nil {
		log.Printf("Error opening file: %v", err)
		return "", err
	}
	defer fileContent.Close()

	// Upload file
	objectName := file.Filename
	contentType := "application/octet-stream"
	_, err = minioClient.PutObject(context.Background(), bucketName, objectName, fileContent, file.Size, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Printf("Error uploading file: %v", err)
		return "", err
	}

	// Generate file URL
	fileURL := fmt.Sprintf("%s/%s/%s", h.cfg.MINIO_ENDPOINT_PUBLIC, bucketName, objectName)

	return fileURL, nil
}

// func (h *helper) AddImageToS3(file *multipart.FileHeader) (string, error) {

// 	// cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-south-1"))
// 	// if err != nil {
// 	// 	fmt.Println("configuration error:", err)
// 	// 	return "", err
// 	// }

// 	// client := s3.NewFromConfig(cfg)
// 	client := s3.NewFromConfig(
// 		aws.Config{Region: "us-east-1",
// 			Credentials: credentials.NewStaticCredentialsProvider(h.cfg.AWS_ACCESS_KEY_ID, h.cfg.AWS_SECRET_ACCESS_KEY, ""),
// 		},
// 		func(o *s3.Options) {
// 			// o.BaseEndpoint = aws.String(h.cfg.AWS_HOST)
// 			o.BaseEndpoint = aws.String("http://ahava.com.vn")
// 		})

// 	uploader := manager.NewUploader(client)

// 	f, openErr := file.Open()
// 	if openErr != nil {
// 		fmt.Println("opening error:", openErr)
// 		return "", openErr
// 	}
// 	defer f.Close()

// 	result, uploadErr := uploader.Upload(context.TODO(), &s3.PutObjectInput{
// 		Bucket: aws.String("ahava"),
// 		Key:    aws.String(file.Filename),
// 		Body:   f,
// 		ACL:    "public-read",
// 	})

// 	if uploadErr != nil {
// 		fmt.Println("uploading error:", uploadErr)
// 		return "", uploadErr
// 	}

// 	return result.Location, nil
// }

func (h *helper) TwilioSetup(username string, password string) {
	client = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: username,
		Password: password,
	})

}

func (h *helper) TwilioSendOTP(phone string, serviceID string) (string, error) {
	to := "+91" + phone
	params := &twilioApi.CreateVerificationParams{}
	params.SetTo(to)
	params.SetChannel("sms")

	resp, err := client.VerifyV2.CreateVerification(serviceID, params)
	if err != nil {

		return " ", err
	}

	return *resp.Sid, nil

}

func (h *helper) TwilioVerifyOTP(serviceID string, code string, phone string) error {

	params := &twilioApi.CreateVerificationCheckParams{}
	params.SetTo("+91" + phone)
	params.SetCode(code)
	resp, err := client.VerifyV2.CreateVerificationCheck(serviceID, params)

	if err != nil {
		return err
	}

	if *resp.Status == "approved" {
		return nil
	}

	return models.ErrValidateOTP
}

func (h *helper) GenerateTokenClients(user models.UserDetailsResponse) (string, error) {
	claims := &AuthCustomClaims{
		ID:    user.ID,
		Email: user.Email,
		Role:  "client",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("ahava"))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (h *helper) GenerateRefferalCode() (string, error) {
	// Calculate the required number of random bytes
	byteLength := (5 * 5) / 8

	// Generate a random byte array
	randomBytes := make([]byte, byteLength)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// Encode the random bytes to base32
	encoder := base32.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZ234567").WithPadding(base32.NoPadding)
	encoded := encoder.EncodeToString(randomBytes)

	// Trim any additional characters to match the desired length
	encoded = encoded[:5]

	return encoded, nil
}

func (h *helper) PasswordHashing(password string) (string, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", models.ErrInternalServer
	}

	hash := string(hashedPassword)
	return hash, nil
}

func (h *helper) CompareHashAndPassword(a string, b string) error {
	err := bcrypt.CompareHashAndPassword([]byte(a), []byte(b))
	if err != nil {
		return err
	}
	return nil
}

func (h *helper) Copy(a *models.UserDetailsResponse, b *models.UserSignInResponse) (models.UserDetailsResponse, error) {
	err := copier.Copy(a, b)
	if err != nil {
		return models.UserDetailsResponse{}, err
	}

	return *a, nil
}
