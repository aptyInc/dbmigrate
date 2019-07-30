package cmd

import (
	"fmt" 
	"github.com/aptyInc/dbmigrate/migrator"
	"github.com/aptyInc/dbmigrate/source"
	"github.com/aptyInc/dbmigrate/target"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
  "os"

	"github.com/spf13/viper"
)


var cfgFile string
var migration *migrator.DBMigrationImplementation

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dbmigrate",
	Short: "Used to migrate database objectcobras",
	Long:  `Used to migrate database objects`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is dbmigrate.yaml in current directory)")
	rootCmd.PersistentFlags().StringP("databaseURL", "u", "postgres://postgres:postgres@localhost:5432/dbmigrate", "the database to which migration has to happen (default is postgres://postgres:postgres@localhost:5432/dbmigrate)")
	rootCmd.PersistentFlags().StringP("directory", "d", ".", "working directory where migration scripts are place")
	viper.BindPFlag("DATABASE_URL", rootCmd.PersistentFlags().Lookup("databaseURL"))
	viper.BindPFlag("directory", rootCmd.PersistentFlags().Lookup("directory"))

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName("dbmigrate")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}



	tgt,err2:=target.GetDatabase(viper.GetString("DATABASE_URL"))
	if err2!= nil {
		fmt.Println(err2)
		os.Exit(1)
	}
	path:=viper.GetString("directory")
	if path == "."{
		home, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		path = home
	}
	src,err1 := source.GetFileSource(path,&source.ReaderImplementation{Fs:afero.NewOsFs()})
	if err1!= nil {
		fmt.Println(err1)
		os.Exit(1)
	}
	migration =  &migrator.DBMigrationImplementation{
		Src:src,
		Tgt:tgt,
	}
}

