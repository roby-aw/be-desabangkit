package user

import (
	userBusiness "api-desatanggap/business/user"
	"api-desatanggap/utils"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"cloud.google.com/go/storage"
	"github.com/labstack/echo/v4"
	"google.golang.org/api/option"
)

type Controller struct {
	service userBusiness.Service
}

func NewController(service userBusiness.Service) *Controller {
	return &Controller{
		service: service,
	}
}

func (Controller *Controller) RegisterAccount(c echo.Context) error {
	var wg sync.WaitGroup
	wg.Add(1)
	Data := userBusiness.RegAccount{}
	c.Bind(&Data)
	var err error
	go func() {
		defer wg.Done()
		_, err = Controller.service.CreateAccount(&Data)
	}()
	wg.Wait()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":     400,
			"messages": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":     200,
		"messages": "success register account",
	})
}

func (Controller *Controller) LoginAccount(c echo.Context) error {
	Data := userBusiness.AuthLogin{}
	c.Bind(&Data)
	result, err := Controller.service.LoginAccount(&Data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":     400,
			"messages": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":     200,
		"messages": "success login",
		"result":   result,
	})
}

// func (Controller *Controller) Registercustomer(c echo.Context) error {
// 	Data := userBusiness.Regcustomer{}
// 	c.Bind(&Data)
// 	result, err := Controller.service.Createcustomer(&Data)
// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, map[string]interface{}{
// 			"code":     400,
// 			"messages": err.Error(),
// 		})
// 	}
// 	return c.JSON(http.StatusOK, map[string]interface{}{
// 		"code":     200,
// 		"messages": "success create data",
// 		"data":     result,
// 	})
// }

func (Controller *Controller) Findcustomer(c echo.Context) error {
	result, err := Controller.service.Findcustomer()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":     400,
			"messages": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":     200,
		"messages": "success get all data customer",
		"data":     result,
	})
}

func (Controller *Controller) UploadPhoto(c echo.Context) error {
	name := c.FormValue("name")
	email := c.FormValue("email")
	file, err := c.FormFile("file")
	err = utils.Upload(name, email, file)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":     400,
			"messages": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":     200,
		"messages": "success upload",
	})
}
func (Controller *Controller) GetRole(c echo.Context) error {
	var wg sync.WaitGroup
	wg.Add(1)
	var err error
	var result []*userBusiness.Role
	go func() {
		defer wg.Done()
		result, err = Controller.service.GetRole()
	}()
	wg.Wait()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":     400,
			"messages": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":     200,
		"messages": "success get role",
		"result":   result,
	})
}

func (Controller *Controller) SmtpEmail(c echo.Context) error {
	email := c.QueryParam("email")
	utils.InitEmail(email)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":     200,
		"messages": "success send email",
		"result":   email,
	})
}

func (Controller *Controller) VerificationAccount(c echo.Context) error {
	code := c.QueryParam("code")
	err := Controller.service.VerificationAccount(code)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":     400,
			"messages": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":     200,
		"messages": "success verification account",
	})
}

func (Controller *Controller) DeleteUser(c echo.Context) error {
	email := c.QueryParam("email")
	err := Controller.service.DeleteUser(email)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":     400,
			"messages": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":     200,
		"messages": "success delete user",
	})
}

func (Controller *Controller) SendVerification(c echo.Context) error {
	email := c.QueryParam("email")
	err := Controller.service.SendVerification(email)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":     400,
			"messages": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":     200,
		"messages": "success send verification",
	})
}

var storageClient storage.Client

func (Controller *Controller) UploadFileHandle(c echo.Context) error {
	f, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})

	}

	blobFile, err := f.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})

	}
	ranstr := utils.RandomString(10)
	f.Filename = *ranstr + ".png"

	err = Controller.UploadFile(blobFile, f.Filename)
	photo_url := "https://storage.googleapis.com/desabangkit-bucket/" + f.Filename
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})

	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":     200,
		"messages": "success upload",
		"result":   photo_url,
	})
}

// UploadFile uploads an object
func (Controller *Controller) UploadFile(file multipart.File, object string) error {
	key_google := os.Getenv("Key_google")
	c, err := storage.NewClient(context.Background(), option.WithCredentialsJSON([]byte(key_google)))
	if err != nil {
		return err
	}
	ctx := context.Background()
	bucket := "desabangkit-bucket"

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	wc := c.Bucket(bucket).Object(object).NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}

	return nil
}

// ctx := appengine.NewContext(c.Request())

// storageClient, err := storage.NewClient(ctx, option.WithCredentialsFile("key.json"))
// if err != nil {
// 	return c.JSON(http.StatusInternalServerError, map[string]interface{}{
// 		"message": "Something went wrong",
// 		"error":   err.Error(),
// 	})
// }
// file, err := c.FormFile("file")
// fmt.Println(file.Filename)
// if err != nil {
// 	return c.JSON(http.StatusInternalServerError, map[string]interface{}{
// 		"message": "Something went wrong",
// 		"error":   err.Error(),
// 	})
// }

// sw := storageClient.Bucket(bucket).Object(file.Filename).NewWriter(ctx)
// // sw.
// u, err := url.Parse("/" + bucket + "/" + sw.Attrs().Name)
// fmt.Println(u.String())
// if err != nil {
// 	return c.JSON(http.StatusInternalServerError, map[string]interface{}{
// 		"message": "Something went wrong",
// 		"error":   err.Error(),
// 	})
// }
// return c.JSON(http.StatusOK, map[string]interface{}{
// 	"message": "File uploaded successfully",
// 	"url":     sw,
// })

func (Controller *Controller) CreateProduct(c echo.Context) error {
	preorder := c.QueryParam("preorder")
	product := userBusiness.InputProduct{}
	c.Bind(&product)
	err := Controller.service.InputProduct(&product, preorder)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":     400,
			"messages": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":     200,
		"messages": "success create product",
	})
}

func (Controller *Controller) GetProductByIdAccount(c echo.Context) error {
	id := c.QueryParam("id")
	product, err := Controller.service.GetProductByIdAccount(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":     400,
			"messages": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":     200,
		"messages": "success get product",
		"result":   product,
	})
}

func (Controller *Controller) GetProductByIdAccStatus(c echo.Context) error {
	id := c.QueryParam("id")
	approved := c.QueryParam("approved")
	var appr, ver bool
	if approved != "" {
		appr, _ = strconv.ParseBool(approved)
	}
	preorder := c.QueryParam("preorder")
	if preorder != "" {
		ver, _ = strconv.ParseBool(preorder)
	}
	product, err := Controller.service.GetProductByIdAccStatus(&appr, &ver, id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":     400,
			"messages": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":     200,
		"messages": "success get product",
		"result":   product,
	})
}

func (Controller *Controller) PostProductTransaction(c echo.Context) error {
	product := userBusiness.InputProductTransaction{}
	c.Bind(&product)
	err := Controller.service.InsertProductTransaction(&product)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":     400,
			"messages": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":     200,
		"messages": "success transaction product",
	})
}

func (Controller *Controller) GetProductTransactionByIDUser(c echo.Context) error {
	id := c.QueryParam("id")
	product, err := Controller.service.GetProductTranscationByIDUser(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":     400,
			"messages": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":     200,
		"messages": "success get product",
		"result":   product,
	})
}
