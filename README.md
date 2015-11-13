goderp
======

goderp abstracts loading configuration files for your app.

Shamelessly borrowing ideas from [derpconf](https://github.com/globocom/derpconf).

rational
--------
goderp provides a way to read simple configurations to your app.

The idea is that you can define variables and default values for them.

After that you'll be able to parse a config file to override the defaults.

At last if environment variables are enabled they are used as the current value.

Usage
-----

    package main
    import (
        "github.com/renatoaquino/goderp"
        "fmt"
    )

    func main() {
        gd := goderp.New()
        gd.Define("PORT", 8888, "Service Port", "Daemon")
        gd.Define("LOG_LEVEL", "info", "Log Level", "Logging")
        
        // Enables override by environment variables 
        gd.EnableEnv()

        gd.Parse("path/to/config.conf")

        port := c.GetInt("PORT")
        log_level := c.GetString("LOG_LEVEL")

        // Dumps the config 
        gd.Dump()        
    }

