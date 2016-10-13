# Bad Gopher

## A study on ransomware, encryption and Go.

Mainly written to learn Go, it's also a proof of concept of how easy is to code something that's quite dangerous. It also sparked my curiosity after reading about the terrible "Fsociety" ransomware, in regards to it's coding practices and design (or complete lack thereof). You may read about it here: https://blog.fortinet.com/2016/09/01/take-it-easy-and-say-hi-to-this-new-python-ransomware

Me, armed with my near-nil knowledge of encryption, went ahead and started writing this thingy. I guess it's interesting how easy is to make a simple ransomware. Just give this thing a filesystem path, and it'll generate an encryption key and shit all over your filesystem encrypting things. And rendering them useless, unless decrypted. That process is around ~70 lines of Go code (FS walking, targetting only some file extensions, and the encryption function).

In 70 lines, someone who knows shit about encryption and a little programming, can fuck up your whole disk (given that you download and execute the binary, which is left as an exercise of social engineering and the like).

As of now (October 13, 2016), the screen presented to the victim is not implemented. The intention is to control the malware locally via a webserver and HTTP requests. A stub of this is coded in `web.go`. Instead, for now, it's controlled via the `main.go` file.  
There are known bugs, and there may probably be more of them, so *use at your own risk*. Which should be needless to say, considering that you are fucking around with malware.


### Disclaimer
Made with the sole purpose of learning, I take absolutely no resposibility on the misuse of this software. As "safe" as it's created to be, it's still dangerous. Use with caution. Or just don't.
