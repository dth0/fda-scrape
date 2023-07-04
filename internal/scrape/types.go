package scrape

import "fmt"

type Resultados map[string]string

type ACAO struct {
	Valor          map[string]string  `json:"valor,omitempty"`
	Ocilacao       map[string]string  `json:"ocilacao,omitempty"`
	Indicadores    map[string]string  `json:"indicadores,omitempty"`
	Demonstrativos map[int]Resultados `json:"resultados,omitempty"`
	Balanco        map[string]string  `json:"balanco,omitempty"`
	Name           string             `json:"name"`
}

func NewACAO(name string) *ACAO {
	return &ACAO{
		Name:           name,
		Valor:          make(map[string]string),
		Ocilacao:       make(map[string]string),
		Indicadores:    make(map[string]string),
		Demonstrativos: make(map[int]Resultados),
		Balanco:        make(map[string]string),
	}
}

type Fis struct {
	Valor       map[string]string  `json:"valor,omitempty"`
	Indicadores map[string]string  `json:"indicadores,omitempty"`
	Balanco     map[string]string  `json:"balanco,omitempty"`
	Imoveis     map[string]string  `json:"imoveis,omitempty"`
	Resultados  map[int]Resultados `json:"resultados,omitempty"`
	Name        string             `json:"name"`
}

func NewFis(name string) *Fis {
	return &Fis{
		Name:        name,
		Valor:       make(map[string]string),
		Indicadores: make(map[string]string),
		Balanco:     make(map[string]string),
		Imoveis:     make(map[string]string),
		Resultados:  make(map[int]Resultados),
	}
}

type Config struct {
	CacheDir   string
	TargerAddr string
	Port       int
	Address    string
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Bind() string {
	return fmt.Sprintf("%s:%d", c.Address, c.Port)
}
