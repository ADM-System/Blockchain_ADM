# Blockchain in Go


## Introduzione

Questo è un esempio di una **blockchain** scritta in **Go**. La blockchain è composta da una sequenza di blocchi, ognuno dei quali contiene:
- Un **indice** univoco,
- Un **timestamp** di quando il blocco è stato creato,
- I **dati** che vengono memorizzati nel blocco,
- L'**hash** del blocco precedente, che lega i blocchi insieme,
- Il **proprio hash** che rappresenta l'identificatore unico di ciascun blocco.

Ogni blocco è creato e collegato alla catena in modo sicuro.

## Funzionalità
1)
- Creazione di blocchi con timestamp.
- Calcolo dell'hash di ogni blocco per garantire l'integrità della blockchain.
- Ogni blocco fa riferimento all'hash del blocco precedente.
  
![image](https://github.com/user-attachments/assets/a9d27b96-55ff-4928-8a65-8b2254126983)

2) aggiunta del Nounce per dare un numero ai blocchi , ed aggiunto un menu' dal quale possiamo aggiungere un blocco , visualizzare l'intera blockchain everificarne la validità e poi uscire.

![image](https://github.com/user-attachments/assets/6ddf2d95-0de5-4a6b-bb4b-dfec5e121c0d)

3) aggiunta un interfaccia grafica web grazsie alla aggiunta si un server web locale scritto ed implementato nel codice in go

 ![image](https://github.com/user-attachments/assets/e1359ac1-92b6-4923-adef-cfa560bb0da8)

 4) aggiunto un sistema array  che tiene tyraccia delle transazioni all'interno del blocco come da foto
  
  ![image](https://github.com/user-attachments/assets/a3ea5470-fe3d-4dfc-ad15-d67d4126dbf7)

5) aggiunta del meccanismo di mining e aggiunto un bottone per il controllo della validità della blockcahin

![image](https://github.com/user-attachments/assets/ab748ee8-f099-433d-9695-4f45ed08fe0b)

6)aggiunto il meccanismo mempool per memorizzare non solo 1 trasnazione ma molte di piu' in un singolo blocco e poi tramite il tasto eseguire il mining del blocco.

7) aggiunto il meccanismo di ricompensa per il miner che viene specificato nel campo , e la sua ricompensa viene inclusa nel blocco.
8) aggiunto calcolo del saldo corrente per ogni utente.
9) aggiunto anche calcolo del saldo per opgni utente con blocco della transazione se il salvo è insufficiente , aggiunta firma e il calcolo del saldo dell'utente , inoltre ho aggiunto un utente admin che parte già con dei coin per test.
    ![image](https://github.com/user-attachments/assets/b5ca7907-be52-4074-ae6e-13508665daa4)
10) funzionalità nuove ma non TESTATE !!!

Modifiche e funzionalità implementate
1. Mempool e mining realistico
Le transazioni vengono prima inserite in un mempool (transazioni in sospeso).
Il mining prende tutte le transazioni dal mempool e le inserisce in un nuovo blocco.
2. Ricompensa per il miner
Ogni blocco minato contiene una transazione di ricompensa (coinbase) per il miner.
3. Calcolo saldo per ogni utente
Puoi vedere il saldo di ogni indirizzo tramite endpoint e interfaccia web.
4. Storico delle transazioni per utente
Puoi vedere tutte le transazioni in cui un indirizzo è coinvolto.
5. Firma digitale reale (ECDSA)
Le transazioni sono firmate digitalmente e la firma viene verificata dal backend.
Ogni transazione contiene anche la chiave pubblica del mittente.
6. Nonce univoco per ogni transazione
Ogni transazione di uno stesso mittente deve avere un nonce crescente, per evitare replay e doppie spese.
7. Propagazione automatica tra nodi (P2P)
I nodi si sincronizzano automaticamente tra loro, adottando la catena più lunga.
Le transazioni vengono propagate tra i peer tramite endpoint dedicato.
8. Gestione dei fork
Quando un nodo adotta una nuova catena, le transazioni che non sono più presenti nei blocchi vengono rimesse nel mempool (recupero delle transazioni orfane).
9. Explorer avanzato
Pagina web dedicata per esplorare la blockchain, cercare blocchi e transazioni, vedere statistiche e dettagli.
10. Documentazione API REST
Pagina HTML che spiega tutti gli endpoint disponibili, con esempi di richiesta e risposta.
11. Miglioramenti di sicurezza
Limite sulla frequenza di mining per utente (anti-spam).
Controlli su importi troppo piccoli o troppo grandi.
Controlli anti-overflow e anti-input malevoli.
Timestamp nelle transazioni e rifiuto di transazioni troppo vecchie o future.
12. Gestione delle transazioni scadute
Le transazioni troppo vecchie (oltre 1 ora) vengono automaticamente rimosse dal mempool, sia periodicamente che prima del mining o dell’aggiunta di una nuova transazione.Modifiche e funzionalità implementate
1. Mempool e mining realistico
Le transazioni vengono prima inserite in un mempool (transazioni in sospeso).
Il mining prende tutte le transazioni dal mempool e le inserisce in un nuovo blocco.
2. Ricompensa per il miner
Ogni blocco minato contiene una transazione di ricompensa (coinbase) per il miner.
3. Calcolo saldo per ogni utente
Puoi vedere il saldo di ogni indirizzo tramite endpoint e interfaccia web.
4. Storico delle transazioni per utente
Puoi vedere tutte le transazioni in cui un indirizzo è coinvolto.
5. Firma digitale reale (ECDSA)
Le transazioni sono firmate digitalmente e la firma viene verificata dal backend.
Ogni transazione contiene anche la chiave pubblica del mittente.
6. Nonce univoco per ogni transazione
Ogni transazione di uno stesso mittente deve avere un nonce crescente, per evitare replay e doppie spese.
7. Propagazione automatica tra nodi (P2P)
I nodi si sincronizzano automaticamente tra loro, adottando la catena più lunga.
Le transazioni vengono propagate tra i peer tramite endpoint dedicato.
8. Gestione dei fork
Quando un nodo adotta una nuova catena, le transazioni che non sono più presenti nei blocchi vengono rimesse nel mempool (recupero delle transazioni orfane).
9. Explorer avanzato
Pagina web dedicata per esplorare la blockchain, cercare blocchi e transazioni, vedere statistiche e dettagli.
10. Documentazione API REST
Pagina HTML che spiega tutti gli endpoint disponibili, con esempi di richiesta e risposta.
11. Miglioramenti di sicurezza
Limite sulla frequenza di mining per utente (anti-spam).
Controlli su importi troppo piccoli o troppo grandi.
Controlli anti-overflow e anti-input malevoli.
Timestamp nelle transazioni e rifiuto di transazioni troppo vecchie o future.
12. Gestione delle transazioni scadute
Le transazioni troppo vecchie (oltre 1 ora) vengono automaticamente rimosse dal mempool, sia periodicamente che prima del mining o dell’aggiunta di una nuova transazione.




  
## Tecnologie utilizzate

- **Go (Golang)**: linguaggio di programmazione utilizzato per la creazione della blockchain.

## Prerequisiti

Per eseguire il progetto, assicurati di avere **Go** installato sul tuo computer. Puoi scaricarlo e installarlo dal sito ufficiale: [https://golang.org/dl/](https://golang.org/dl/)
