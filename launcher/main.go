//go:build windows

package main

import (
 "os"
 "os/exec"
 "path/filepath"
 "syscall"
 "unsafe"
 "context"
 "time"
 "github.com/LuciferSoria/Minerux-Apocalypse/launcher/service"
)

const (
 wsOverlappedWindow = 0x00CF0000
 wsVisible = 0x10000000
 wsChild = 0x40000000
 wmDestroy = 0x0002
 wmCommand = 0x0111
 swShow = 5
 idPlay = 101
 idStatus = 102
 idOfficial = 103
)
var (
 user32 = syscall.NewLazyDLL("user32.dll")
 kernel32 = syscall.NewLazyDLL("kernel32.dll")
 shell32 = syscall.NewLazyDLL("shell32.dll")
 registerClass = user32.NewProc("RegisterClassW")
 createWindow = user32.NewProc("CreateWindowExW")
 defWindowProc = user32.NewProc("DefWindowProcW")
 showWindow = user32.NewProc("ShowWindow")
 updateWindow = user32.NewProc("UpdateWindow")
 getMessage = user32.NewProc("GetMessageW")
 translateMessage = user32.NewProc("TranslateMessage")
 dispatchMessage = user32.NewProc("DispatchMessageW")
 postQuitMessage = user32.NewProc("PostQuitMessage")
 messageBox = user32.NewProc("MessageBoxW")
 getModuleHandle = kernel32.NewProc("GetModuleHandleW")
 shellExecute = shell32.NewProc("ShellExecuteW")
)
type wndClass struct {
 Style uint32
 WndProc uintptr
 ClsExtra int32
 WndExtra int32
 Instance uintptr
 Icon uintptr
 Cursor uintptr
 Background uintptr
 MenuName *uint16
 ClassName *uint16
}
type point struct { X, Y int32 }
type msg struct {
 Hwnd uintptr
 Message uint32
 WParam uintptr
 LParam uintptr
 Time uint32
 Pt point
 Private uint32
}
func utf16(s string) *uint16 { return syscall.StringToUTF16Ptr(s) }
func alert(hwnd uintptr, text string) {
 messageBox.Call(hwnd, uintptr(unsafe.Pointer(utf16(text))), uintptr(unsafe.Pointer(utf16("Minerux Apocalypse"))), 0x40)
}
func openOfficialLauncher(hwnd uintptr) {
 for _, base := range []string{os.Getenv("ProgramFiles(x86)"), os.Getenv("ProgramFiles")} {
  exe := filepath.Join(base, "Minecraft Launcher", "MinecraftLauncher.exe")
  if _, err := os.Stat(exe); err == nil {
   if err := exec.Command(exe).Start(); err == nil { return }
  }
 }
 result, _, _ := shellExecute.Call(hwnd, uintptr(unsafe.Pointer(utf16("open"))), uintptr(unsafe.Pointer(utf16("minecraft:"))), 0, 0, swShow)
 if result <= 32 {
  alert(hwnd, "No se ha encontrado Minecraft Launcher. Instálalo desde la web oficial. Esta compilación DEV todavía no instala Minecraft ni NeoForge.")
 }
}
func checkStatus(hwnd uintptr) {
 endpoint:=os.Getenv("MINERUX_API_URL")
 if endpoint=="" { alert(hwnd,"Backend no configurado. Define MINERUX_API_URL con la URL HTTPS del servicio beta.");return }
 client,err:=service.New(endpoint)
 if err!=nil {alert(hwnd,"Configuracion de backend invalida: "+err.Error());return}
 ctx,cancel:=context.WithTimeout(context.Background(),15*time.Second)
 defer cancel()
 status,err:=client.GetStatus(ctx)
 if err!=nil {alert(hwnd,"No se puede consultar el servidor: "+err.Error());return}
 switch status.Status {
 case "maintenance": alert(hwnd,"MANTENIMIENTO: "+status.Message)
 case "beta": alert(hwnd,"BETA PRIVADA: Necesitas una invitacion valida para jugar. "+status.Message)
 case "online": alert(hwnd,"SERVIDOR DISPONIBLE: "+status.Message)
 }
}
// launchAfterStatusCheck fails closed: maintenance, private beta and unreachable API block play.
// Full beta authorization and direct game launch are not implemented yet.
func launchAfterStatusCheck(hwnd uintptr) {
 endpoint:=os.Getenv("MINERUX_API_URL")
 if endpoint=="" {alert(hwnd,"No hay backend configurado. JUGAR no esta disponible todavia.");return}
 client,err:=service.New(endpoint)
 if err!=nil {alert(hwnd,"URL del backend invalida.");return}
 ctx,cancel:=context.WithTimeout(context.Background(),15*time.Second)
 defer cancel()
 s,err:=client.GetStatus(ctx)
 if err!=nil {alert(hwnd,"No se puede verificar el estado del servidor. Se ha bloqueado JUGAR.");return}
 if s.Status=="maintenance" {alert(hwnd,"Servidor en mantenimiento: "+s.Message);return}
 if s.Status=="beta" {alert(hwnd,"Beta privada: el inicio de sesion y la verificacion de invitaciones siguen pendientes.");return}
 alert(hwnd,"El servidor indica disponible, pero el lanzamiento directo de Minecraft aun no esta implementado. Puedes usar ABRIR MINECRAFT OFICIAL.")
}
func windowProc(hwnd uintptr, message uint32, wparam, lparam uintptr) uintptr {
 switch message {
 case wmCommand:
  if uint16(wparam & 0xffff) == idPlay { go launchAfterStatusCheck(hwnd); return 0 }
  if uint16(wparam & 0xffff) == idStatus { go checkStatus(hwnd); return 0 }
  if uint16(wparam & 0xffff) == idOfficial { openOfficialLauncher(hwnd); return 0 }
 case wmDestroy:
  postQuitMessage.Call(0)
  return 0
 }
 r, _, _ := defWindowProc.Call(hwnd, uintptr(message), wparam, lparam)
 return r
}
func child(hwnd uintptr, class, title string, x,y,w,h int32, id uintptr) {
 createWindow.Call(0, uintptr(unsafe.Pointer(utf16(class))), uintptr(unsafe.Pointer(utf16(title))), uintptr(wsChild|wsVisible), uintptr(x), uintptr(y), uintptr(w), uintptr(h), hwnd, id, 0, 0)
}
func main() {
 instance, _, _ := getModuleHandle.Call(0)
 name := utf16("MineruxApocalypseWindow")
 wc := wndClass{WndProc: syscall.NewCallback(windowProc), Instance: instance, Background: 6, ClassName: name}
 if r, _, _ := registerClass.Call(uintptr(unsafe.Pointer(&wc))); r == 0 { panic("RegisterClassW failed") }
 hwnd, _, _ := createWindow.Call(0, uintptr(unsafe.Pointer(name)), uintptr(unsafe.Pointer(utf16("Minerux Apocalypse - DEV"))), wsOverlappedWindow|wsVisible, 160, 120, 640, 360, 0, 0, instance, 0)
 if hwnd == 0 { panic("CreateWindowExW failed") }
 child(hwnd, "STATIC", "MINERUX APOCALYPSE - THE FALLEN WORLD", 35, 38, 550, 35, 0)
 child(hwnd, "STATIC", "COMPILACION DE DESARROLLO - NO ES EL LAUNCHER FINAL", 35, 87, 550, 30, 0)
 child(hwnd, "STATIC", "Pendiente: Microsoft OAuth, NeoForge, mods y entrada directa al servidor.", 35, 125, 550, 45, 0)
 child(hwnd, "BUTTON", "JUGAR (EN DESARROLLO)", 35, 185, 350, 45, idPlay)
 child(hwnd, "BUTTON", "CONSULTAR ESTADO", 395, 185, 195, 45, idStatus)
 child(hwnd, "BUTTON", "ABRIR MINECRAFT OFICIAL", 35, 240, 350, 40, idOfficial)
 showWindow.Call(hwnd, swShow)
 updateWindow.Call(hwnd)
 var m msg
 for {
  r, _, _ := getMessage.Call(uintptr(unsafe.Pointer(&m)), 0, 0, 0)
  if int32(r) <= 0 { break }
  translateMessage.Call(uintptr(unsafe.Pointer(&m)))
  dispatchMessage.Call(uintptr(unsafe.Pointer(&m)))
 }
}
