package middleware_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/nvasilev98/ginkgo-gomega/pkg/middleware"
	"github.com/nvasilev98/ginkgo-gomega/pkg/middleware/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("JWT Parser", func() {
	const tokenWithoutHeader = "Bearer payload.signature"

	var (
		mockCtrl     *gomock.Controller
		mockContext  *gin.Context
		recorder     *httptest.ResponseRecorder
		jwtParser    *middleware.JWTParser
		jwtValidator *mocks.MockJWTValidator
	)

	BeforeEach(func() {
		recorder = httptest.NewRecorder()
		mockContext, _ = gin.CreateTestContext(recorder)
		mockContext.Request = httptest.NewRequest("GET", "http://test.com", nil)
		mockCtrl = gomock.NewController(GinkgoT())
		jwtValidator = mocks.NewMockJWTValidator(mockCtrl)
		jwtParser = middleware.NewJWTParserWithOptions(middleware.WithValidator(jwtValidator))
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	When("token extraction fails", func() {
		It("should return status code unauthorized", func() {
			jwtParser.Handle(mockContext)
			Expect(recorder.Code).To(Equal(http.StatusUnauthorized))
			var errResp string
			Expect(json.Unmarshal(recorder.Body.Bytes(), &errResp)).To(Succeed())
			Expect(errResp).To(Equal("Authentication failed"))
		})
	})

	When("token validation fails", func() {
		BeforeEach(func() {
			mockContext.Request.Header.Set("Authorization", tokenWithoutHeader)
			jwtValidator.EXPECT().ValidateJWT(gomock.Any(), gomock.Any()).Return(errors.New("err"))
		})

		It("should return status code unauthorized", func() {
			jwtParser.Handle(mockContext)
			Expect(recorder.Code).To(Equal(http.StatusUnauthorized))
			var errResp string
			Expect(json.Unmarshal(recorder.Body.Bytes(), &errResp)).To(Succeed())
			Expect(errResp).To(Equal("Invalid JWT token"))
		})
	})
})
