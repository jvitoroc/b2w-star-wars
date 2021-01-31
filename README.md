
# b2w-star-wars

API desenvolvida para o desafio técnico da vaga de Desenvolvedor Go na B2W. 

## Testes

Para testar a aplicação, executar o seguinte comando dentro da pasta raiz do projeto:

> go test -v

## Documentação

### [POST] Criar um novo planeta
> hostname:port/planet/

**Exemplo de corpo**
```json
{
    "name": "Tatooine",
    "climate": "Arid",
    "terrain": "Dessert"
}
```
**Exemplo de resposta**
```json
{
    "message": "The planet was successfully created.",
    "planet": {
        "id": "6015b48eccd6e8fa2e01f4d8",
        "name": "Tatooine",
        "climate": "Arid",
        "terrain": "Dessert",
        "filmsAppearedIn": 5
    }
}
```
___
### [GET] Listar planetas
> hostname:port/planet/

**Exemplo de resposta**
```json
{
    "message": "The planets were successfully retrieved.",
    "planets": [
        {
            "id": "6015b54bccd6e8fa2e01f4db",
            "name": "Tatooine",
            "climate": "arid",
            "terrain": "desert",
            "filmsAppearedIn": 5
        },
        {
            "id": "6015b565ccd6e8fa2e01f4dc",
            "name": "Alderaan",
            "climate": "temperate",
            "terrain": "grasslands, mountains",
            "filmsAppearedIn": 2
        }
    ]
}
```
___
### [GET] Buscar por ID
> hostname:port/planet/{id}

**Exemplo de URL:** hostname:port/planet/6015b565ccd6e8fa2e01f4dc

**Exemplo de resposta**
```json
{
    "message": "The planet was successfully retrieved.",
    "planet": {
        "id": "6015b565ccd6e8fa2e01f4dc",
        "name": "Alderaan",
        "climate": "temperate",
        "terrain": "grasslands, mountains",
        "filmsAppearedIn": 2
    }
}
```
___
### [GET] Buscar por nome
> hostname:port/planet/?search={criteria}

**Exemplo de URL:** hostname:port/planet/?search=Aldera

**Exemplo de resposta**
```json
{
    "message": "The planets were successfully retrieved.",
    "results": [
        {
            "id": "6015b565ccd6e8fa2e01f4dc",
            "name": "Alderaan",
            "climate": "temperate",
            "terrain": "grasslands, mountains",
            "filmsAppearedIn": 2
        }
    ]
}
```
___
### [DELETE] Remover um planeta
> hostname:port/planet/{id}

**Exemplo de URL:** hostname:port/planet/6015859defc67f00159d44a3

**Exemplo de resposta**
```json
{
    "message": "The planet was successfully deleted."
}
```
