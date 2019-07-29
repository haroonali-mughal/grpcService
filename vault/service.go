package vault

import(
	"context"
	"golang.org/x/crypto/bcrypt"
        "github.com/go-kit/kit/endpoint"
	"encoding/json"
	"log"
	"net/http"
        //"vault/endpoint"
	"errors"
)

type Service interface {
	Hash(ctx context.Context , password []string) ([]string , error)
	Validate(ctx context.Context , password , hash []string) (bool , error)
}

type vaultService struct {}

func NewService() Service {
	log.Println("inside vault.NewService()")

	return vaultService{}
}

func (vaultService) Hash(ctx context.Context,password []string) ([]string , error) {
	//arrayLen := len(password)
	var hashArray []string

	for i := 0 ; i<2 ; i++  {
		hash,err := bcrypt.GenerateFromPassword([]byte(password[i]) , bcrypt.DefaultCost)
		hashArray = append(hashArray,string(hash))
		if err != nil {
			return nil,err
		}
	}
	return hashArray[:] , nil


	/*hash,err := bcrypt.GenerateFromPassword([]byte(password) , bcrypt.DefaultCost)
	if err != nil {
		return "",err
	}
	return string(hash) , nil*/
}

func (vaultService) Validate(ctx context.Context , password , hash []string) (bool , error) {
	for i := 0 ; i<len(password) ; i++ {
		err := bcrypt.CompareHashAndPassword([]byte(hash[i]),[]byte(password[i]))
		if err != nil {
			return false , nil
		}
	}
	return true , nil

	/*err := bcrypt.CompareHashAndPassword([]byte(hash),[]byte(password))
	if err != nil {
		return false , nil
	}
	return true, nil*/
}

type hashRequest struct {
	Password []string `json:"password"`
}

type hashResponse struct {
	Hash []string `json :"password"`
	Err string `json : "err,omitempty"`
}

func decodeHashRequest(ctx context.Context , r *http.Request) (interface{},error){
	var req hashRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil , err
	}
	return req , nil
}

type validateRequest struct {
	Password []string `json:"password"`
	Hash []string `json:"hash"`
}

type validateResponse struct {
	Valid bool `json:"valid"`
	Err string `json:"err,omitempty"`
}

func decodeValidateRequest(ctx context.Context,r *http.Request) (interface{},error){
	var req validateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil , err
	}
	return req , nil

}

func encodeResponse(ctx context.Context , w http.ResponseWriter , response interface{}) error{
	return json.NewEncoder(w).Encode(response)
}

func MakeHashEndpoint(srv Service) endpoint.Endpoint {
	log.Println("isnide MakeHashEndpoint in service.go")
	return func(ctx context.Context , request interface{}) (interface{},error){
		req := request.(hashRequest)
		v, err := srv.Hash(ctx , req.Password[:])
		if err != nil {

			log.Println("thre is some error in MakeHashEndpoint")
			return hashResponse{v,err.Error()}, nil
		}
		return hashResponse{v,""},nil
	}	
}

func MakeValidateEndpoint(srv Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
        req := request.(validateRequest)
        v, err := srv.Validate(ctx, req.Password[:], req.Hash[:])
        if err != nil {
        	return validateResponse{false, err.Error()}, nil
        }
        return validateResponse{v, ""}, nil
	}
}

type Endpoints struct {
	HashEndpoint endpoint.Endpoint
	ValidateEndpoint endpoint.Endpoint
}

func (e Endpoints) Validate(ctx context.Context, password, hash []string) (bool, error) {
	req := validateRequest{Password: password, Hash: hash}
	resp, err := e.ValidateEndpoint(ctx, req)
	if err != nil {
		return false, err
	}
	validateResp := resp.(validateResponse)
	if validateResp.Err != "" {
		return false, errors.New(validateResp.Err)
	}
	return validateResp.Valid, nil
}

func (e Endpoints) Hash(ctx context.Context , password []string) ([]string , error){
	req := hashRequest{Password : password}
	resp, err := e.HashEndpoint(ctx , req)
	if err != nil {
		return nil,err
	}
	hashResp := resp.(hashResponse)
	if hashResp.Err != "" {
		return nil,errors.New(hashResp.Err)
	}
	return hashResp.Hash, nil
}
