package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("frontend/*.html"))
}

type Pelanggan struct {
	IDPelanggan string
	NoTelp      string
	Email       string
	Nama        string
	NoRumah     string
	NamaJalan   string
	KodePos     string
}

type SemuaPelanggan struct {
	AllPelanggan []Pelanggan
}

type Penjual struct {
	IDPenjual string
	Email     string
	Nama      string
	NoTelp    string
	Alamat    string
}

type SemuaPenjual struct {
	AllPenjual []Penjual
}

type Pembayaran struct {
	IDPembayaran    string
	IDPelanggan     string
	IDPesanan       string
	JenisPembayaran string
}

type SemuaPembayaran struct {
	AllPembayaran []Pembayaran
}

type Produk struct {
	IDProduk  string
	IDPenjual string
	Nama      string
	Jenis     string
	Stok      int
	Harga     int
}

type SemuaProduk struct {
	Product []Produk
}

func returnPelanggan(emailUser string) []Pelanggan {
	var pelanggan Pelanggan
	var arrPelanggan []Pelanggan

	db := connect()
	rows, err := db.Query("SELECT * FROM Pelanggan WHERE email = ?", emailUser)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&pelanggan.IDPelanggan,
			&pelanggan.NoTelp,
			&pelanggan.Email,
			&pelanggan.Nama,
			&pelanggan.NoRumah,
			&pelanggan.NamaJalan,
			&pelanggan.KodePos,
		)
		if err != nil {
			log.Fatal(err)
		}
		arrPelanggan = append(arrPelanggan, pelanggan)
	}
	return arrPelanggan
}

func returnPenjual(emailSeller string) string {
	var penjual Penjual
	var dataPenjual []Penjual

	db := connect()
	rows, err := db.Query("SELECT * FROM Penjual WHERE email = ?", emailSeller)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer rows.Close()

	var returnable string
	for rows.Next() {
		err := rows.Scan(
			&penjual.IDPenjual,
			&penjual.Email,
			&penjual.Nama,
			&penjual.NoTelp,
			&penjual.Alamat,
		)
		if err != nil {
			log.Fatal(err)
		}
		dataPenjual = append(dataPenjual, penjual)
		returnable = penjual.Nama
	}

	if len(dataPenjual) > 0 {
		return "Welcome, " + returnable
	}
	return "Tidak ada email yang terdaftar"
}

func returnProduk(yangDicari string) []Produk {
	var produk Produk
	var arrProduk []Produk

	yangDicari = "%" + yangDicari + "%"

	db := connect()
	rows, err := db.Query("SELECT * FROM Produk WHERE nama LIKE ?", yangDicari)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&produk.IDProduk,
			&produk.IDPenjual,
			&produk.Nama,
			&produk.Jenis,
			&produk.Stok,
			&produk.Harga,
		)
		if err != nil {
			log.Fatal(err)
		}
		arrProduk = append(arrProduk, produk)
	}
	return arrProduk
}

func returnPayment(dicari string) []Pembayaran {
	var pembayaran Pembayaran
	var dataPembayaran []Pembayaran

	db := connect()
	rows, err := db.Query("SELECT * FROM Pembayaran WHERE ID_Pembayaran = ?", dicari)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&pembayaran.IDPembayaran,
			&pembayaran.IDPelanggan,
			&pembayaran.IDPesanan,
			&pembayaran.JenisPembayaran,
		)
		if err != nil {
			log.Fatal(err)
		}
		dataPembayaran = append(dataPembayaran, pembayaran)
	}
	return dataPembayaran
}

//Section ini Kebawah Fungsi Endpoint

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		tpl.ExecuteTemplate(w, "index.html", nil)
	}

}

func loginCustHandler(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "loginCustomer.html", nil)
}

func loginCustCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/custlogin", http.StatusSeeOther)
		return
	}

	email := r.FormValue("emailuser")

	returnedAuth := returnPelanggan(email)

	returnedData := SemuaPelanggan{
		returnedAuth,
	}

	tpl.ExecuteTemplate(w, "customerCheck.html", returnedData)
}

func loginSellHandler(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "loginSeller.html", nil)
}

func loginSellCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/sellerlogin", http.StatusSeeOther)
		return
	}

	email := r.FormValue("emailseller")

	sellEmail := returnPenjual(email)

	returnedEmail := struct {
		SellerEmail string
	}{
		SellerEmail: sellEmail,
	}

	tpl.ExecuteTemplate(w, "sellerCheck.html", returnedEmail)
}

func productHandler(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "productPage.html", nil)
}

func SearchProd(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/shopping", http.StatusSeeOther)
		return
	}

	searchValue := r.FormValue("prod")

	AllProd := returnProduk(searchValue)

	allProduct := SemuaProduk{
		AllProd,
	}

	tpl.ExecuteTemplate(w, "productShow.html", allProduct)
}

func cekPembayaran(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "paymentCheck.html", nil)
}

func searchPembayaran(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/payment", http.StatusSeeOther)
		return
	}

	searchPay := r.FormValue("idbayar")

	AllPayment := returnPayment(searchPay)

	AllPay := SemuaPembayaran{
		AllPayment,
	}

	tpl.ExecuteTemplate(w, "paymentStatus.html", AllPay)
}

func registerLanding(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "register.html", nil)
}

func registerCustomer(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	id := r.FormValue("id")
	notelp := r.FormValue("notelp")
	email := r.FormValue("email")
	nama := r.FormValue("nama")
	noRumah := r.FormValue("norumah")
	namaJalan := r.FormValue("jalan")
	kodepos := r.FormValue("kodepos")

	db := connect()
	insert, err := db.Query("INSERT INTO Pelanggan VALUES(?,?,?,?,?,?,?)", id, notelp, email, nama, noRumah, namaJalan, kodepos)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer insert.Close()

	tpl.ExecuteTemplate(w, "regSuccess.html", nil)
}

func registerSellerLanding(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "registerSeller.html", nil)
}

func registerSeller(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/registerseller", http.StatusSeeOther)
		return
	}

	id := r.FormValue("id")
	email := r.FormValue("email")
	name := r.FormValue("nama")
	noTelp := r.FormValue("telepon")
	alamat := r.FormValue("alamat")

	db := connect()
	insert, err := db.Query("INSERT INTO Penjual VALUES(?,?,?,?,?)", id, email, name, noTelp, alamat)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer insert.Close()

	tpl.ExecuteTemplate(w, "regSuccess.html", nil)
}

func insertProduct(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "insertProduct.html", nil)
}

func productQuerying(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/sellproduct", http.StatusSeeOther)
		return
	}

	id := r.FormValue("id")
	idPenjual := r.FormValue("idpenjual")
	nama := r.FormValue("nama")
	jenis := r.FormValue("jenis")
	stock := r.FormValue("stock")
	harga := r.FormValue("harga")

	db := connect()
	insert, err := db.Query("INSERT INTO Produk VALUES(?,?,?,?,?,?)", id, idPenjual, nama, jenis, stock, harga)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer insert.Close()

	tpl.ExecuteTemplate(w, "prodSuccess.html", nil)
}

func beliBarang(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "buyProduct.html", nil)
}

func beliBarangTrans(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/buyproduct", http.StatusSeeOther)
		return
	}

	idProd := r.FormValue("id")
	idCust := r.FormValue("idc")
	idSell := r.FormValue("ids")
	idPes := r.FormValue("idp")
	quan := r.FormValue("quan")

	db := connect()
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		fmt.Println(err.Error())
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO Pesanan VALUES(?,?,?,'2020-12-11',?,(SELECT harga FROM Produk WHERE ID_Produk = ?),5000,(5000+?*(SELECT harga FROM Produk WHERE ID_Produk = ?)))", idPes, idCust, idSell, quan, idProd, quan, idProd)
	if err != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		return
	}

	_, err = tx.ExecContext(ctx, "UPDATE Produk SET stok = stok - ? WHERE ID_Produk = ?", quan, idProd)
	if err != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		return
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO Produk_Pesanan VALUES(?,?)", idProd, idPes)
	if err != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Succesfully Commited")
	}

	tpl.ExecuteTemplate(w, "reservationSuccess.html", nil)
}

func bayarBarangHandler(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "paymentTry.html", nil)
}

func prosesPembayaran(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/pay", http.StatusSeeOther)
		return
	}

	idpemb := r.FormValue("id")
	idcust := r.FormValue("idpelanggan")
	idpes := r.FormValue("idpesanan")
	noRek := r.FormValue("norek")
	metode := r.FormValue("metode")

	db := connect()
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		fmt.Println(err.Error())
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO Pembayaran VALUES(?, ?, ?, 'Daring')", idpemb, idcust, idpes)
	if err != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		return
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO PembayaranDaring VALUES(?, ?, ?)", idpemb, noRek, metode)
	if err != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Succesfully Commited")
	}

	tpl.ExecuteTemplate(w, "onlinePaymentSuccess.html", nil)
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/register", registerLanding)
	http.HandleFunc("/regconfirm", registerCustomer)
	http.HandleFunc("/registerseller", registerSellerLanding)
	http.HandleFunc("/regsellconfirm", registerSeller)
	http.HandleFunc("/sellproduct", insertProduct)
	http.HandleFunc("/sellcheck", productQuerying)
	http.HandleFunc("/buyproduct", beliBarang)
	http.HandleFunc("/buyprodcheck", beliBarangTrans)
	http.HandleFunc("/custlogin", loginCustHandler)
	http.HandleFunc("/custcheck", loginCustCheck)
	http.HandleFunc("/sellerlogin", loginSellHandler)
	http.HandleFunc("/sellercheck", loginSellCheck)
	http.HandleFunc("/shopping", productHandler)
	http.HandleFunc("/prodsearch", SearchProd)
	http.HandleFunc("/payment", cekPembayaran)
	http.HandleFunc("/paymentsearch", searchPembayaran)
	http.HandleFunc("/pay", bayarBarangHandler)
	http.HandleFunc("/payconfirm", prosesPembayaran)
	http.ListenAndServe(":14045", nil)
}
