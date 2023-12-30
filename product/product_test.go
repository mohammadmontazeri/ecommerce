package product

// import (
// 	"ecommerce/models"
// 	"errors"
// 	"testing"

// 	// "github.com/elliotchance/redismock/v8"
// 	"github.com/go-redis/redis/v8"
// 	"github.com/go-redis/redismock/v8"
// 	"github.com/stretchr/testify/assert"
// 	// "github.com/stretchr/testify/require"
// )

// // func TestUploadImage(t *testing.T) {

// // 	body := new(bytes.Buffer)
// // 	writer := multipart.NewWriter(body)
// // 	// create a new form-data header name data and filename data.txt
// // 	dataPart, err := writer.CreateFormFile("file", "6.jpg")
// // 	require.NoError(t, err)

// // 	// copy file content into multipart section dataPart
// // 	f, err := os.Open("6.jpg")
// // 	require.NoError(t, err)
// // 	_, err = io.Copy(dataPart, f)
// // 	require.NoError(t, err)
// // 	require.NoError(t, writer.Close())

// // 	// create HTTP request & response
// // 	r, err := http.NewRequest(http.MethodPost, "/files", body)
// // 	require.NoError(t, err)
// // 	r.Header.Set("Content-Type", writer.FormDataContentType())
// // 	w := httptest.NewRecorder()
// // 	///
// // 	gin.SetMode(gin.TestMode)

// // 	ctx, _ := gin.CreateTestContext(w)
// // 	// ctx.Request = &http.Request{
// // 	// 	Header: make(http.Header),
// // 	// }
// // 	// ctx.Request.Method = "GET"
// // 	// ctx.Request.Header.Set("Content-Type", "multipart/form-data")
// // 	// ctx.Writer.Header().Set("Content-Type", "multipart/form-data")

// // 	_, err = UploadProductImage(ctx)

// // 	assert.Nil(t, err)

// // }

// var (
// 	client *redis.Client
// )

// func TestRedisSetProduct(t *testing.T) {
// 	// mr, err := miniredis.Run()
// 	// if err != nil {
// 	// 	log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	// }

// 	// client = redis.NewClient(&redis.Options{
// 	// 	Addr: mr.Addr(),
// 	// })

// 	// m := redismock.NewNiceMock(client)

// 	// m.On("HSet", "key", mock.AnythingOfType("models.product")).Return(redis.NewStatusResult("", nil))

// 	// r := NewRedisRepository(m)

// 	// product := models.Product{}
// 	// err = r.SetProduct("key", product)
// 	// assert.NoError(t, err)
// 	db, mock := redismock.NewClientMock()

// 	key := "key"
// 	product := models.Product{
		
// 	}
// 	mock.ExpectHSet(key, product).SetErr(errors.New("FAIL"))
// 	repo := NewRedisRepository(db)

	
// 	err := repo.SetProduct(key, product)

// 	assert.NoError(t, err)

// }
