package service

import (
 "bytes"
 "context"
 "encoding/json"
 "errors"
 "fmt"
 "io"
 "net/http"
 "net/url"
 "strings"
 "time"
)
type Status struct { Status string `json:"status"`; Message string `json:"message"` }
type Access struct { Access bool `json:"access"` }
type Client struct { base string; http *http.Client }
func New(base string) (*Client,error) {
 u,err:=url.Parse(base)
 if err!=nil || u.Scheme!="https" || u.Host=="" || u.User!=nil || u.RawQuery!="" || u.Fragment!="" || (u.Path!="" && u.Path!="/") || u.Opaque!="" {return nil,errors.New("backend URL must be HTTPS without credentials, query or fragment")}
 return &Client{base:strings.TrimRight(u.String(),"/"),http:&http.Client{Timeout:12*time.Second,CheckRedirect:func(req *http.Request,via []*http.Request) error {if len(via)>=1{return errors.New("backend redirects are not allowed")};return nil}}},nil
}
func (c *Client) request(ctx context.Context,method,path,token string,body any,out any) error {
 var reader io.Reader
 if body!=nil {data,err:=json.Marshal(body);if err!=nil{return err};reader=bytes.NewReader(data)}
 req,err:=http.NewRequestWithContext(ctx,method,c.base+path,reader);if err!=nil{return err}
 if token!=""{req.Header.Set("Authorization","Bearer "+token)}
 if body!=nil{req.Header.Set("Content-Type","application/json")}
 resp,err:=c.http.Do(req);if err!=nil{return err};defer resp.Body.Close()
 if resp.StatusCode<200 || resp.StatusCode>=300 {return fmt.Errorf("backend returned HTTP %d",resp.StatusCode)}
 return json.NewDecoder(io.LimitReader(resp.Body,1<<20)).Decode(out)
}
func (c *Client) GetStatus(ctx context.Context) (Status,error) {
 var s Status
 err:=c.request(ctx,http.MethodGet,"/api/v1/status","",nil,&s)
 if err!=nil{return s,err}
 switch s.Status {case "online","maintenance","beta":return s,nil;default:return s,errors.New("unrecognized server status")}
}
func (c *Client) CheckAccess(ctx context.Context,playerToken string) (bool,error) {
 if playerToken==""{return false,errors.New("verified player token required")}
 var a Access
 err:=c.request(ctx,http.MethodGet,"/api/v1/me/access",playerToken,nil,&a)
 return a.Access,err
}
func (c *Client) Redeem(ctx context.Context,playerToken,code string) (bool,error) {
 if playerToken==""{return false,errors.New("verified player token required")}
 if len(code)<16 || len(code)>128 {return false,errors.New("invalid beta code length")}
 var a Access
 err:=c.request(ctx,http.MethodPost,"/api/v1/beta/redeem",playerToken,map[string]string{"code":code},&a)
 return a.Access,err
}
func (c *Client) CanPlay(ctx context.Context,playerToken string) (bool,string,error) {
 s,err:=c.GetStatus(ctx)
 if err!=nil{return false,"No se puede comprobar el estado del servidor",err}
 if s.Status=="maintenance"{return false,s.Message,nil}
 if s.Status=="online"{return true,"",nil}
 allowed,err:=c.CheckAccess(ctx,playerToken)
 if err!=nil{return false,"No se puede verificar el acceso beta",err}
 if !allowed{return false,"Necesitas un código beta válido",nil}
 return true,"",nil
}
