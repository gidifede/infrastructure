# Infrastructure Repository - Documentazione

## Panoramica

La soluzione "Logistic-backbone" è un framework di gestione eventi basato sull'architettura event-sourcing focalizzato sulla logistica. Informazioni architetturali di dettaglio sono reperibili qui [inserire link]. La presente documentazione è rivolta alla descrizione degli aspetti implementativi e non architetturalie nonchè alla comprensione della struttura del seguente repo ed alle scelte di dettaglio.
Il presente repo è progettato per il deployment su AWS di un ambiente robusto e scalabile che possa gestire vari aspetti come la configurazione, la logica di business, l'interfaccia API e i test.

### Scopo

L'obiettivo è automatizzare e semplificare il processo di deploy e gestione di servizi cloud, consentendo una rapida iterazione e scalabilità. La soluzione copre:

- Definizione di schemi e configurazioni
- Gestione di servizi attraverso handler di comandi e query
- Documentazione e test

## Architettura del repo

### Configuration

- **Ruolo**: Definire schemi e configurazioni in formato JSON.
- **Tecnologie**: JSON Schema

### Infrastructure

- **Ruolo**: Contiene il codice per il deploy delle risorse cloud.
- **Tecnologie**: AWS CDK, Python
- **Interazione**: Fa il deploy dei servizi definiti nella directory `services`.

### Services

- **Ruolo**: Contiene la logica di business e il codice sorgente per vari servizi.
- **Tecnologie**: Linguaggio di programmazione specifico al servizio (es. Python, Go)
- **Interazione**: Viene deployato attraverso il codice nella directory `infrastructure`.

### Docs

- **Ruolo**: Fornire la documentazione del progetto.
- **Tecnologie**: Markdown, altri formati di documentazione

### Swagger

- **Ruolo**: Documentazione delle API.
- **Tecnologie**: Swagger

### Tests

- **Ruolo**: Contiene test unitari e funzionali.
- **Tecnologie**: Framework di test specifici al linguaggio (es. pytest per Python)

## Flusso di Lavoro

1. **Definizione della Configurazione**: Gli schemi e le configurazioni vengono definiti nella directory `configuration`.
2. **Sviluppo del Servizio**: La logica di business e il codice sorgente vengono sviluppati nella directory `services`.
3. **Deploy**: Il codice nella directory `infrastructure` fa il deploy delle risorse cloud e dei servizi.
4. **Documentazione e Test**: La documentazione è mantenuta aggiornata e i test vengono eseguiti per garantire la qualità del codice.

## Esempio di Deploy con CDK

```python
from aws_cdk import (
    aws_lambda as _lambda,
    aws_apigateway as apigw,
    core
)

class MyStack(core.Stack):

    def __init__(self, scope: core.Construct, id: str, **kwargs) -> None:
        super().__init__(scope, id, **kwargs)

        # Definizione della funzione Lambda
        my_lambda = _lambda.Function(
            self, 'MyLambda',
            runtime=_lambda.Runtime.PYTHON_3_8,
            handler='my_lambda.handler',
            code=_lambda.Code.from_asset('services/my_service')
        )

        # API Gateway
        apigw.LambdaRestApi(
            self, 'MyApi',
            handler=my_lambda,
        )
```

# Get started

Si suggerisce l'utilizzo di Taskfile:
https://taskfile.dev/installation/

## Comandi utili di taskfile

di seguito alcuni comandi utili (descrizione riportata in lingua inglese)
| Nome Comando                                                   | Descrizione                                         |
| --------------------------------------------------------------| ----------------------------------------------------|
| `task init`                                                    | Initialize python dev environment                   |
| `task sync-services`                                           | Update all submodule repos with remote              |
| `task services-switch-branch -- [master|develop|branchname]`  | Switch to the specified branch all submodule repos  |
| `task update-graph`                                            | Update the architecture image based on current CDK script |
| `task run-unit-tests`                                          | Update all services and run only unit tests         |
| `task run-integration-tests` or `task run-integration-tests -- [aws-profile]` | Update all services and run only integration tests  |
