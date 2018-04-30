package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"regexp"
	"strconv"
	"time"

	"github.com/flibustier/multichain-client"
)

///////////////// S T R U C T U R E S ///////////////////

//AddressBalance is used in order to return the adresses and their ammount of money.
type AddressBalance struct {
	Addresses   string
	Balances    float64
	CurrentTime string
}

func main() {
	heure := time.Now().Local()
	fmt.Println("Il est : ", heure.Format("2006-01-02T15:04:05.999999999Z07:00"))

	user, erra := user.Current()
	if erra != nil {
		panic(erra)
	}
	path := user.HomeDir + "\\AppData\\Roaming\\Multichain"
	err12 := os.Chdir(path)
	if err12 != nil {
		fmt.Println(err12)
	}
	cmd := exec.Command(".\\multichaind.exe", "Amacoin")
	err31 := cmd.Start()
	if err31 != nil {
		fmt.Printf("Est-ce que ça tourne ?: %s \n ", err31)
	}

	duration := time.Duration(3) * time.Second
	time.Sleep(duration) // little sleep  (3s) before connecting
	////////////////////////// Démarrage de multichaind.exe Amacoin@IP:Port
	////////////////////////// Pour se connecter au noeud Papa (pas besoin de Ip:Port si on est le noeud papa)

	// Connexion to the holly blockchain hosting the noble écu
	// We need a central node, used as a DNS seed
	///////////////////////// FLAGS TO LAUNCH THE .EXE WITH OPTIONS ////////////////////////////
	chain := flag.String("chain", "Amacoin", "is the name of the chain")
	host := flag.String("host", "localhost", "is a string for the hostname")
	port := flag.Int("port", 4336, "is a number for the host port")
	username := flag.String("username", "multichainrpc", "is a string for the username")
	password := flag.String("password", "DYiL6vb71Y8qfEo9CkYr5wyZ3GqjRxrjzkYyjsA9S1k2", "is a string for the password")
	flag.Parse()

	logs := GetLogins(*chain)
	*username = logs[0]
	*password = logs[1]
	*port = GetPort(*chain)

	////////////////////////
	client := multichain.NewClient(
		*chain,
		*username,
		*password,
		*port,
	).ViaNode(
		*host,
		*port,
	)

	///////////////////////////WHILE CONNECTED TO A CHAIN /////////////////////////////
	obj, err := client.GetInfo() // Returns general information about this node and blockchain.
	if err != nil {
		fmt.Println("Il y a une erreur dans la connexion de la chaine. ")
		panic(err)
	}
	//fmt.Println(obj)

	obj, err = client.GetAddresses(false) // Get the addresses in our wallet.
	if err != nil {                       // Impossible to reach our wallet, please ask for lost objects.
		log.Fatal("[FATAL] Could not get addresses from Multichain", err)
	}

	addresses := obj.Result().([]interface{}) // Different addresses stored on the node
	//fmt.Println(addresses)                                     // Array with the addresses
	//log.Printf("[OK] Your main address is %s\n", addresses[0]) // First wallet

	///////////////////////////////// Asset Definition ////////////////////////////////
	RewardName := "Amacoin" // Nom de notre monnaie.
	InitialReward := 10.0   // Récompense d'entrée.
	cents := 0.01           // Unité monétaire divisionnaire de l'écu.
	///////////////////////////////////////////////////////////////////////////////////

	address := addresses[0].(string)                                         // The first wallet is the principle one. End of discussion
	obj, err = client.Issue(true, address, RewardName, InitialReward, cents) // If it's the first time the node is launched, we have to create the asset for reward

	if err != nil { // Asset already existing
		log.Printf("[OK] Asset %s seemsto be already existing", RewardName)
	} else { // Creation of the non existing asset
		log.Printf("[OK] Mon adresse est toujours : %s", address)
		obj, err = client.IssueMore(address, RewardName, 10) // Noob award ?
		if err != nil {
			log.Printf("[ERREUR SUR L'ADRESSE]")
		} else {
			log.Printf("[OK] ON A RAJOUTE L'ARGENT") // Award granted
		}
		log.Printf("[OK] Asset %s successfuly created", RewardName) // Graphical confirmation of the asset creation's success
	}
	// End of the initialization of da wallet.
	log.Printf("[OK] On a bien démarré notre noeud. La bourse est disponible à l'adresse : %s", address)

	fmt.Printf("\n [OK] ========================================== \n \n \n")
	res, err5 := GetWalletBalances(client)
	if err5 != nil {
		fmt.Printf("%s \n", err5)
		fmt.Printf("résultat : \n %s \n", res)
	}

	//////////////////////////////////////////////////////////

	address1 := Identification(client)
	_, erreu := choice(client, address1, RewardName)
	fmt.Printf("%s \n %s \n", address1, erreu)
}

//////////////////////////////////// A C T I O N S /////////////////////////////////

//	Function that will call sub functions depending on the user's choice.
func choice(client *multichain.Client, address string, RewardName string) (int, error) {
	c := exec.Command("clear") // Efface l'écran
	c.Stdout = os.Stdout
	c.Run()
	var res int
	for true {
		c := exec.Command("clear") // Efface l'écran
		c.Stdout = os.Stdout
		c.Run()

		fmt.Printf("============== MENU ============== \n Choisir une action parmi les actions proposees suivantes \n Et entrer le numéro correspondant à la suite. \n Quitter 0) \n Menu Administrateur 1) \n Menu Client 2) \n ============================== \n")
		_, err := fmt.Scanf("%d\n", &res)
		if err != nil { // SCAN is Not OK
			fmt.Printf("Wrong imput, please try again.\n")
			//return 0, err
		} else { // Scan is OK
			fmt.Printf("On a choisi : %d ", res)
			switch res {
			case 1: // Admin menu
				fmt.Printf("Menu Admin \n")
				err = ChoiceAdmin(client, RewardName)
			case 2: // Client menu
				fmt.Printf("Menu Client \n")
				err = ChoiceClient(client, RewardName, address)
			case 0: // Exit
				//fmt.Println("Exiting.")
				//os.Exit(0)
				return 0, nil
			default:
				fmt.Println("Not an option")
			}
			//return 0, err
		}
	}
	return 0, nil
}

//Identification is a function that asks very basically the user to inform the program his office
func Identification(client *multichain.Client) string {
	var res int
	tableau := GetLocalAddresses(client)
	fmt.Printf("\n ============ I D E N T I F I C A T I O N ============= \n")
	fmt.Printf("Les adresses disponibles sur le noeud sont: \n")
	for i := range tableau {
		fmt.Printf("Adresse %d: %s \n", i, tableau[i])
	}
	fmt.Printf("======================================================= \n Quelle adresse correspond à votre bureau? Entrer le numéro correspondant.\n")
	_, err := fmt.Scanf("%d\n", &res)
	if err != nil { // SCAN is Not OK
		fmt.Printf("Wrong imput, please try again.\n")
		return ""
	}
	res1 := tableau[res]
	return res1
}

//SendTransaction is a function that will send to the server the data about the last transaction
func SendTransaction(client *multichain.Client, address string) error {
	url := "http://localhost:8989/transaction"
	fmt.Println("URL:>", url)
	resp1, erreur := GetAddressTransaction(client, address)
	if erreur != nil {
		fmt.Printf("Error in the SendTransaction \n %s \n", erreur)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(resp1))
	//fmt.Printf("Transaction : \n %s \n \n %s \n", resp1, req)
	//req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")
	client1 := &http.Client{}
	resp, err := client1.Do(req)
	if err != nil {
		b := errors.New("Impossible de communiquer avec le serveur les informations de transaction")
		return b
	}
	defer resp.Body.Close()
	fmt.Println("response Status:", resp.Status)
	return nil
}

//SendMoney is a function that allows one to send assets to another address
func SendMoney(client *multichain.Client, res3 string, asset string) error {
	var res int
	var qt float64

	tableau := GetGlobalAdresses(client)
	fmt.Printf("______________________________\nLes adresses disponibles sont: \n")
	for i := range tableau {
		fmt.Printf("Adresse %d: %s \n", i, tableau[i])
	}
	fmt.Printf("============================ \n Quelle adresse créditer? Entrer le numéro correspondant.\n")
	_, err := fmt.Scanf("%d\n", &res)
	if err != nil { // SCAN is Not OK
		fmt.Printf("Wrong imput, please try again.\n")
		return err
	}
	res1 := tableau[res]

	fmt.Printf("Quelle quantité d'argent transférer ?\n")
	_, err2 := fmt.Scanf("%f\n", &qt)
	if err2 != nil { // SCAN is Not OK
		fmt.Printf("Wrong imput, please try again.\n Erreur:%s \n", err2)
		return err2
	}
	_, err3 := client.SendAssetFrom(res3, res1, asset, qt) // Send "qty" of "asset" from "res3" to "res1". The first parameter returned is a Response Type that contain the transaction Id that can be useful

	duration := time.Duration(10) * time.Second
	time.Sleep(duration)

	if err3 != nil {
		fmt.Printf("Cannot send Asset. \n")
		return err3
	}
	err4 := SendTransaction(client, res3)
	if err4 != nil {
		fmt.Printf("Cannot transmit transaction.")
	}
	return nil // Everything is all right.
}

// ChoiceAdmin is a function that open the Menu for admin functions.
func ChoiceAdmin(client *multichain.Client, asset string) error {
	c := exec.Command("clear") // Efface l'écran
	c.Stdout = os.Stdout
	c.Run()
	var res1 int
	fmt.Printf("=========== MENU ADMIN ========== \n Creer une nouvelle adresse dans le portefeuille 1) \n Crediter une adresse 2) \n Sortie 0) \n ============================== \n")
	_, err := fmt.Scanf("%d\n", &res1)
	switch res1 {
	case 1: // Transfer assets
		err := CreateAddress(client)
		if err != true {
			fmt.Printf("Error in CreateAddress")
			return nil
		}
	case 2: // Issue asset
		err := IssueMoney(client, asset)
		if err != true {
			fmt.Printf("Error in IssueMoney")
			return nil
		}
	case 0: // Exit
		//fmt.Println("Exiting...")
		//os.Exit(0)
		return nil
	default:
		fmt.Println("Not an option")
	}
	//fmt.Printf("J'ai rentré %d, il y a erreur : %s \n", res1, err)
	if err != nil { // SCAN is Not OK
		fmt.Printf("Wrong imput, please try again.\n")
		return err
	}

	return nil
}

// ChoiceClient is a Function that starts a menu for the Client Options
func ChoiceClient(client *multichain.Client, RewardName string, address string) error {
	c := exec.Command("clear") // Efface l'écran
	c.Stdout = os.Stdout
	c.Run()
	var res1 int
	fmt.Printf("=========== MENU CLIENT ========== \n Consulter son portefeuille 1) \n Consulter les adresses 2) \n Virement 3) \n Sortie 0) \n ============================== \n")
	_, err := fmt.Scanf("%d\n", &res1)
	//fmt.Printf("J'ai rentré %d, il y a erreur : %s \n", res1, err)
	if err != nil { // SCAN is Not OK
		fmt.Printf("Wrong imput, please try again.\n")
		return err
	}
	// Scan is OK
	//fmt.Printf("On a choisi : %d \n", res1)
	switch res1 {
	case 1: // Check our wallet
		_, err := GetWalletBalances(client)
		if err != nil {
			fmt.Printf("Error \n")
		}
	case 2: // Peer addresses
		fmt.Println("_________________________________ \nLes adresses disponibles sont: ")
		GetGlobalAdresses(client)

	case 3: // Transfer assets
		err := SendMoney(client, address, RewardName)
		if err != nil {
			fmt.Printf("Asset can not be sent")
			return err
		}
	case 0: // Exit
		//	fmt.Println("Exiting...")
		//os.Exit(0)
		return nil
	default:
		fmt.Println("Not an option")
	}
	return nil
}

///////////////////////////////////// A C C S E S S E U R S ////////////////////////////////////

// GetGlobalAdresses is a function that returns an array of the available adresses
func GetGlobalAdresses(client *multichain.Client) []string {
	c := exec.Command("clear") // Efface l'écran
	c.Stdout = os.Stdout
	c.Run()
	tabret := make([]string, 0)
	params := []interface{}{"receive"}
	msg := client.Command( // It will do the manual command
		"listpermissions", // listpermissions that returns the allowed to receive a transaction
		params,            // Basically all the addresses of the network
	)
	coucou, erre := client.Post(msg)
	if erre != nil {
		fmt.Printf("Erreur cli post %s \n", erre)
	}

	for j := range coucou.Result().([]interface{}) { // Here we want to extract the addresses
		fmt.Printf(" ===================== \n %d ) ", j) // From the structure in coucou
		plop := coucou.Result().([]interface{})[j].(map[string]interface{})
		plip := plop["address"].(string)
		tabret = append(tabret, plip) // Adding the addresses
		fmt.Printf("%s \n ==================== \n", tabret[j])
	}
	var input string
	fmt.Scanln(&input)
	return tabret
}

//GetAddressTransaction is a function that return the json of the last tranaction involving the address.
func GetAddressTransaction(client *multichain.Client, address string) ([]byte, error) {
	response, a := client.ListAddressTransactions(address, 1, 0, false) // Get the last transaction (1), without skipping(0), and without too mmuch text (verbose = false)
	if a != nil {
		err := errors.New("Error in the function GetAddressTransaction")
		var ret []byte
		return ret, err
	}
	res2B, _ := json.Marshal(response) // transform the response in json response.
	return res2B, nil
}

//GetLocalAddresses is a function that return a list of the addresses contained in da wallet.
func GetLocalAddresses(client *multichain.Client) []string {
	obj, err := client.GetAddresses(false) // Get the addresses in our wallet.
	if err != nil {                        // Impossible to reach our wallet, please ask for lost objects.
		log.Fatal("[FATAL] Could not get addresses from Multichain", err)
	}
	addresses := obj.Result().([]interface{}) // Different addresses stored on the node
	adresses := make([]string, 0)
	for i := range addresses {
		adresses = append(adresses, addresses[i].(string))
	}
	return adresses
}

//GetAddressBalance is a function that return the ammount of money stored in the given address
func GetAddressBalance(client *multichain.Client, address string) float64 {
	money, err := client.GetAddressBalances(address) // To see how much we have in the *address'wallet
	if err != nil {
		log.Printf("Je n'arrive pas à ouvrir la bourse et compter les pièces :'( ")
		//	} else { // Good way to print the wallet's contents.
		//		log.Printf("\n Vous avez actuellement %f Amacoin.\n", money.Result().([]interface{})[0].(map[string]interface{})["qty"].(float64))
	}
	if len(money.Result().([]interface{})) == 0 {
		return 0.0
	}
	ret := money.Result().([]interface{})[0].(map[string]interface{})["qty"].(float64)
	return ret
}

//GetLogins Is a function that will read the multichain.conf file and returns user login and password.
func GetLogins(chain string) []string {
	user, err := user.Current()
	if err != nil {
		log.Fatal("[FATAL] Could not get user from Multichain", err)
	}

	login := "NULL"              // Case in which we cannot find any login.
	password := "NULL"           // Case in which we cannot find any password.
	path1 := user.HomeDir + "\\" //////////////////// PATH DIRECTORY FOR WINDOWS USERS \\\\\\\\\\\\\\\\\\\\\\\\\\\\\\
	path2 := "AppData\\Roaming\\Multichain\\"
	path3 := chain + "\\"
	path4 := "multichain.conf"
	path := path1 + path2 + path3 + path4
	inFile, err1 := os.Open(path)

	if err1 != nil {
		log.Fatal("[FATAL] Could not open Multichain path", err1)
	}

	re := regexp.MustCompile("rpcpassword=([a-zA-Z0-9]+)") // Gonna search for those strings followed by alphanumerics symbols
	re1 := regexp.MustCompile("rpcuser=([a-zA-Z0-9]+)")

	defer inFile.Close()
	scanner := bufio.NewScanner(inFile) // Scan the file
	scanner.Split(bufio.ScanLines)      // Scan by Lines
	tableau := make([]string, 0)        // Tableau will store the data
	for scanner.Scan() {                // We read the file line by line
		//%ùfmt.Println(scanner.Text())
		if re.MatchString(scanner.Text()) { // If the line matches the searched string (after the defined string)
			password = re.FindStringSubmatch(scanner.Text())[1] // Get the scanned text
			//fmt.Println(password)
		} else if re1.MatchString(scanner.Text()) {
			login = re1.FindStringSubmatch(scanner.Text())[1]
			//fmt.Println(login)
		}
	}
	tableau = append(tableau, login)
	tableau = append(tableau, password) // Keep tablea growing with the matched strings
	return tableau
}

//GetPort Is a function that will read the params.dat file and returns the default port.
func GetPort(chain string) int {
	user, err := user.Current() // Get user's name
	if err != nil {
		log.Fatal("[FATAL] Could not get user from Multichain", err)
	}

	port := "NULL"               // Case in which we cannot find any port.
	path1 := user.HomeDir + "\\" // Windows Path in which multichain needs to be installed
	path2 := "AppData\\Roaming\\Multichain\\"
	path3 := chain + "\\"
	path4 := "params.dat"
	path := path1 + path2 + path3 + path4
	inFile, err1 := os.Open(path) // Open path

	if err1 != nil {
		log.Fatal("[FATAL] Could not open Multichain params.dat", err)
	}

	re := regexp.MustCompile("default-rpc-port = ([0-9]+)") //We want to get the number after "default-rpc-port = "

	defer inFile.Close()
	scanner := bufio.NewScanner(inFile) // Scanner file
	scanner.Split(bufio.ScanLines)      // Scan by line

	for scanner.Scan() { //We read the file line by line
		//fmt.Println(scanner.Text())
		if re.MatchString(scanner.Text()) { //If it matches
			port = re.FindStringSubmatch(scanner.Text())[1] //Get the matched text
			//fmt.Println(port)
		}
	}
	port1, err := strconv.Atoi(port) //convert to integers.
	return port1
}

//GetWalletBalances Is a function that will return the summed ammount of all the assets contained in the addresses
func GetWalletBalances(client *multichain.Client) ([]byte, error) {
	c := exec.Command("clear") // Efface l'écran
	c.Stdout = os.Stdout
	c.Run()
	tab := GetLocalAddresses(client)    // Get the wallet addresses
	length := len(tab)                  // Number of addresses stored
	balances := make([]float64, length) // Create array of float, flexible length
	var total float64                   // Total that will store the total money stored in the wallet
	total = 0                           // init
	tabul := make([]AddressBalance, length)
	heure := time.Now().Local()
	timestamp := heure.Format("2006-01-02T15:04:05.999999999Z07:00")
	for i := 0; i < len(tab); i++ {
		//balances[i] = float64(i)*0.212 + 1
		balances[i] = GetAddressBalance(client, tab[i])
		fmt.Printf("L'addresse %d contient %f AmaCoin \n", i, balances[i])
		total = total + balances[i]

		tabul[i] = AddressBalance{ // Put that in a struct
			Addresses:   tab[i],
			Balances:    balances[i],
			CurrentTime: timestamp,
		}

	}
	fmt.Printf("Pour un total de: %f AmaCoin \n", total)
	// Put in mapper amount of money indexed by the address corresponding

	res2B, ck := json.Marshal(tabul) // Convert to JSON
	if ck != nil {
		fmt.Printf("impossible to convert Wallet balances into JSON \n")
		fmt.Printf("%s", ck)
		return nil, errors.New("Can not convert wallet balances into JSON")
	}
	return res2B, nil // Return the result and an error code (here everything is Ok)
}

////////////////////  A D M I N  C O M M A N D S  ////////////////////
//																	//
// grant: connect receive send admin mine							//
//																	//
// revoke: connect receive send admin mine							//
//																	//
// Issue asset to an address										//
//																	//
//////////////////////////////////////////////////////////////////////

//CreateAddress is a function that creates a new address within the wallet and grant them with the basic permissions
func CreateAddress(client *multichain.Client) bool {
	c := exec.Command("clear") // Efface l'écran
	c.Stdout = os.Stdout
	c.Run()
	res, err := client.GetNewAddress()
	if err != nil {
		fmt.Printf("Impossible de créer la nouvelle adresse. \n %s \n", err)
		return false
	}
	permissions := []string{"connect", "send", "receive", "mine"}
	resTr := []string{res.Result().(string)}
	resp, erroer := client.Grant(resTr, permissions)
	if erroer != nil {
		fmt.Printf("Grant denied : \n %s \n", erroer)
	}
	fmt.Printf("Nouvelle adresse créée avec succès. \n %s \n ======================== \n", resp)
	return true
}

//IssueMoney is a function that allows to credit some money to an user choosen address.
func IssueMoney(client *multichain.Client, asset string) bool {
	c := exec.Command("clear") // Efface l'écran
	c.Stdout = os.Stdout
	c.Run()
	var res int
	var qt float64

	tableau := GetGlobalAdresses(client) // Get Addresses
	fmt.Printf("______________________________\nLes adresses disponibles sont: \n")
	for i := range tableau {
		fmt.Printf("Adresse %d: %s \n", i, tableau[i])
	}
	fmt.Printf("============================ \n Quelle adresse créditer? Entrer le numéro correspondant.\n")
	_, err := fmt.Scanf("%d\n", &res)
	if err != nil { // SCAN is Not OK
		fmt.Printf("Wrong imput, please try again.\n")
		return false
	}
	res1 := tableau[res]

	fmt.Printf("Quelle quantité d'argent créer ?\n")
	_, err2 := fmt.Scanf("%f\n", &qt)
	if err2 != nil { // SCAN is Not OK
		fmt.Printf("Wrong imput, please try again.\n Erreur:%s \n", err2)
		return false
	}

	rei, err54 := client.IssueMore(res1, asset, qt)
	if err54 != nil {
		fmt.Printf("Impossible de créer la monnaie sur l'adresse choisie.\n %s \n", rei)
		return false
	}
	return true
}
