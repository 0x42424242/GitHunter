## GitHunter

This program will check the nominated domains for exposed /.git/ directories in the hostname webroot. 

### Arguments
`-filename` 	

The filename which contains domains without http or https prefixes. One per line. 

`-quiet`

Supresses the majority of the output and only prints a . for each assessed domain unless a directory is found. 


### Installation

`go build GitHunter.go` or alternatively use a release binary. 


### Legend
```
[+] - .git directory found
[-] - No .git directory found
[!] - Error attempting to retrieve the domain
```


### Example #1

Check all the domains listed in the files. 

`./GitHunter.exe -filename domains.txt`


#### Output
```
[-] [200] meet.ssp.example.com.au
[+] [200] [HTTPS: true] sip.ssp.example.com.au
[-] [404] ae.example.com.au
[!] extranet.example.com.au
[-] [302] cybersecuritystrategy.example.com.au
[!] webconf.ssp.example.com.au
[-] [200] meet.example.com.au
[-] [404] ea.example.com.au
```


### Example #2

Check all the domains listed in the files and have quiet output

`./GitHunter.exe -filename domains.txt -quiet`

#### Output
```
..........
[+] [200] [HTTPS: true] sip.ssp.example.com.au
.............
```
