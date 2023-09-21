package cmd

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"

	"github.com/fatih/color"
	"github.com/nub06/iis-hero/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var credCmd = &cobra.Command{
	Use:   "cred",
	Short: "Display current credentials",
	Long:  `cred`,
	Run: func(cmd *cobra.Command, args []string) {

		remotehost := viper.GetString("remotehost")
		remoteusername := viper.GetString("remoteusername")
		remotepassword := viper.GetString("remotepassword")
		remotedomain := viper.GetString("remotedomain")

		m := map[string]string{

			"Config Profile": "",
			"Hostname":       remotehost,
			"Username":       remoteusername,
			"Password":       remotepassword,
			"Domain":         remotedomain,
		}

		forcePass, _ := cmd.Flags().GetBool("force")

		for k, v := range m {

			if k == "Config Profile" {

				err, config := service.ShowCurrentConfig()
				if err != nil {

					fmt.Printf("%s: Empty\n", k)
				} else {
					configProfile := config[0].Profile

					if configProfile != "" {

						fmt.Printf("%s: %s \n", k, configProfile)

					}
				}
			}

			if isStrNotEmpty(v) {

				s := decrypt(v)

				if s != "" {
					if k == "Password" {
						if forcePass {
							fmt.Printf("%s: %s \n", k, s)
						} else {
							fmt.Printf("%s: %s \n", k, passReplaceAsterisk(s))
						}

					} else {
						fmt.Printf("%s: %s \n", k, s)
					}

				} else {

					if k != "Config Profile" {

						fmt.Printf("%s: Empty\n", k)
					}
				}

			} else {
				if k != "Config Profile" {
					fmt.Printf("%s: Empty\n", k)
				}
			}

		}
	},
}

func init() {

	credCmd.Flags().BoolP("force", "f", false, "Use this command if you want to see password without asterisk")
}

func setRemoteComputerDetails() service.RemoteComputer {

	c := new(service.RemoteComputer)

	if viper.GetString("remotehost") == "" {

		log.Fatal(color.HiGreenString("Check your credentials. -c, --computer flag can not be empty \ne.g:iis-hero login cred"))
	}

	c.ComputerName = decrypt(viper.GetString("remotehost"))
	c.DomainName = decrypt(viper.GetString("remotedomain"))
	c.Password = decrypt(viper.GetString("remotepassword"))
	c.Username = decrypt(viper.GetString("remoteusername"))
	return *c
}

func clearViperInfo() {

	viper.Set("remotehost", "")
	viper.Set("remoteusername", "")
	viper.Set("remotepassword", "")
	viper.Set("remotedomain", "")
	viper.Set("isRemote", false)
	viper.WriteConfig()

}

func passReplaceAsterisk(password string) string {
	length := len(password)
	if length <= 3 {
		if length == 3 {
			return string([]rune{rune('*'), '*', '*'})
		}
		return password
	}
	runes := []rune(password)
	for i := 1; i < length-1; i++ {
		runes[i] = '*'
	}
	return string(runes)
}

func encrypt(str string) string {

	key := []byte("$E-x2wY*n_h&?2mP")
	plaintext := []byte(str)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	cipherText := make([]byte, aes.BlockSize+len(plaintext))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plaintext)

	return base64.URLEncoding.EncodeToString(cipherText)
}

func decrypt(cryptoText string) string {
	key := []byte("$E-x2wY*n_h&?2mP")

	cipherText, _ := base64.URLEncoding.DecodeString(cryptoText)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	if len(cipherText) < aes.BlockSize {
		fmt.Println(cipherText)
		panic("ciphertext too short")
	}
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText)
}

func isStrNotEmpty(str string) bool {

	return str != ""

}
