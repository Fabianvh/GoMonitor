## Algemeen 
Monitoring tool voor de challenge semester 2 Infrastructure
Hiervoor heb ik de tool gossm gebruikt.
Gossm is een monitoring tool gebouwd door ene Silvio Simunic uit Croatie. Deze man heeft een tool geschreven voor het monitoren van servers.

Omdat ik vond dat er functies in de code zaten die ik niet nodig of zelf wou herschrijven als ik tijd had heb ik deze uit de code gehaald. Dit ging om functies zoals het versturen van berichten over telegram, mail en sms heb ik uit de code gehaald. Ook heb ik wat foutmeldingen uit de code gehaald die erin zijn gekomen door een verouderde versie van go. Op de server waar ik de code op draai heb ik de nieuwste versie van golang geïnstalleerd. Bepaalde parameters moesten worden aangepast omdat de compiler deze niet kon inlezen.



## Instructies
### Go code draaien
Voor het downloaden van de code
`#Go get github.com/Fabianvh/GoMonitor`

Deze zal het configuratiefile `configs/myconfig.json` gebruiken, deze is aan te passen naar onderstaande uitleg over het configuratiebesatnd

Draaien is `#go run .` in de map waar main.go staat

## Configuratiebestand
Het aanpassen van het configuratie bestand gaat op de volgende manier

#### Settings
•	‘Checkinterval’ geeft het aantal secondes aan waarin de check worden uitgevoerd. 

•	‘Timeout’ is het aantal keer dat de server gaat proberen de check uit te voeren voordat hij aan zal geven dat het niet bereikbaar is. 

•	‘maxConnections’ staat voor het maximale aantal connecties dat de monitor goedkeurt. 

•	‘exponentialBackoffSeconds’ zijn het aantal secondes die het programma toevoegt zodra deze een fout tegenkomt. Dus zodra 1 foutmelding optreed doet hij vijf seconde voordat hij het opnieuw probeert. Gebeurd dit nogmaals voegt hij nogmaals vijf seconde toe dus wacht hij tien seconde

#### Servers
Bij het kopje daaronder (Servers) kunnen server toegevoegd kunnen worden die gecheckt moeten worden. 

•	Bij ‘name’ geef je de titel op die op de website weergegeven moet worden.

•	‘IpAddress’ is het ip waarop de server of services bereikbaar is.

•	‘port’ is de poort waarop er gecontroleerd moet worden.

•	‘protocol’ is het protocol waarmee de checks worden uitgevoerd.

•	‘checkInterval’ per server kan je nogmaals aangeven bij welke interval gecontroleerd word.

•	‘timeout’ staat voor het aantal pogingen dat de server zal doen om te controleren of het apparaat beschikbaar is.


```json
{
    "settings": {
        "monitor": {
            "checkInterval": 30,   
            "timeout": 5,
            "maxConnections": 50,
            "exponentialBackoffSeconds": 5
        }
    },
"servers": [
{
"name":"Extern - Google",
"ipAddress":"8.8.8.8",
"port": 80,
"protocol": "tcp",
"checkInterval": 5,
"timeout": 5
},
{
"name":"CH_RT01",
"ipAddress":"10.20.0.1",
"port": 80,
"protocol": "tcp",
"checkInterval": 5,
"timeout": 5
},
{
"name":"CH_AD01",
"ipAddress":"10.10.0.10",
"port": 80,
"protocol": "tcp",
"checkInterval": 5,
"timeout": 5
},
{
"name":"CH_APP01",
"ipAddress":"10.10.0.12",
"port": 80,
"protocol": "tcp",
"checkInterval": 5,
"timeout": 5
},
{
"name":"CH_BP01",
"ipAddress":"10.10.0.11",
"port": 80,
"protocol": "udp",
"checkInterval": 5,
"timeout": 5
}
]
}
```
