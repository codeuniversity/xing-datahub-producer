package handler

import (
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/codeuniversity/xing-datahub-protocol"
	"github.com/golang/protobuf/jsonpb"

	"net/http/httptest"

	"github.com/Shopify/sarama"
	"github.com/Shopify/sarama/mocks"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRequestHandlerForUsers(t *testing.T) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = false
	config.Producer.Return.Errors = false
	producer := mocks.NewAsyncProducer(t, config)
	handler := RequestHandler{
		Producer:   producer,
		RawMessage: &protocol.RawUser{},
		Topic:      "users",
	}
	Convey("Given a correct user", t, func() {
		user := &protocol.User{Id: 1}
		marshaler := &jsonpb.Marshaler{}
		message, err := marshaler.MarshalToString(user.Parse())
		So(err, ShouldBeNil)
		reader := strings.NewReader(message)
		Convey("With no access-token", func() {
			req := httptest.NewRequest("POST", "/users", reader)
			resp := httptest.NewRecorder()
			Convey("When there is a token defined in os.env", func() {
				os.Setenv("token", "123")
				handler.ServeHTTP(resp, req)

				So(resp.Code, ShouldEqual, http.StatusUnauthorized)
			})

			Convey("When there is no token defined in os.env", func() {
				os.Setenv("token", "")
				producer.ExpectInputAndSucceed()
				handler.ServeHTTP(resp, req)

				So(resp.Code, ShouldNotEqual, http.StatusUnauthorized)

			})
		})
		Convey("With the correct access-token", func() {
			req := httptest.NewRequest("POST", "/users", reader)
			req.Header.Add("access-token", "123")
			resp := httptest.NewRecorder()
			Convey("When there is a token defined in os.env", func() {
				os.Setenv("token", "123")
				producer.ExpectInputAndSucceed()
				handler.ServeHTTP(resp, req)

				So(resp.Code, ShouldNotEqual, http.StatusUnauthorized)
			})

			Convey("When there is no token defined in os.env", func() {
				os.Setenv("token", "")
				producer.ExpectInputAndSucceed()
				handler.ServeHTTP(resp, req)

				So(resp.Code, ShouldNotEqual, http.StatusUnauthorized)
			})
			Convey("We answer with the correct code", func() {
				handler.ServeHTTP(resp, req)
				producer.ExpectInputAndSucceed()
				So(resp.Code, ShouldEqual, http.StatusOK)
			})
		})
	})
	Convey("Given an invalid user", t, func() {
		user := &protocol.User{Id: 1}
		marshaler := &jsonpb.Marshaler{}
		message, err := marshaler.MarshalToString(user)
		So(err, ShouldBeNil)
		reader := strings.NewReader(message)
		req := httptest.NewRequest("POST", "/users", reader)
		resp := httptest.NewRecorder()

		os.Setenv("token", "")
		handler.ServeHTTP(resp, req)
		So(resp.Code, ShouldEqual, http.StatusBadRequest)
	})
}
