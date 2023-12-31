package mem_measure

type item struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Quantity  int    `json:"quantity"`
	Unit_cost int    `json:"unit_cost"`
}

var product = []item{
	{ID: 0, Name: "Product 1", Quantity: 10, Unit_cost: 100},
	{ID: 1, Name: "Product 2", Quantity: 20, Unit_cost: 150},
	{ID: 2, Name: "Product 3", Quantity: 100, Unit_cost: 10},
}

type MemoryStorage struct {
	data map[int]item
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data: make(map[int]item),
	}
}

func (s *MemoryStorage) GET(ID int) (Name string, Quantity int, Unit_cost int) {
	for _, p := range product {
		if p.ID == ID {
			return p.Name, p.Quantity, p.Unit_cost
		}
	}
	return
}

func (s *MemoryStorage) GETALL() []item {
	return product
}

func (s *MemoryStorage) POST(Name string, Quantity int, Unit_cost int) (ID int) {
	var NewProduct item
	NewProduct.ID = product[len(product)-1].ID + 1
	NewProduct.Name = Name
	NewProduct.Quantity = Quantity
	NewProduct.Unit_cost = Unit_cost
	product = append(product, NewProduct)
	return NewProduct.ID
}

func (s *MemoryStorage) DELETE(ID int) {
	for i, p := range product {
		if p.ID == ID {
			product = append(product[:i], product[i+1:]...)
			break
		}
	}
}

func (s *MemoryStorage) PUT(ID int, Name string, Quantity int, Unit_cost int) {

	for i, p := range product {
		if p.ID == ID {
			product[i].Name = Name
			product[i].Quantity = Quantity
			product[i].Unit_cost = Unit_cost
		}
	}
}

func (s *MemoryStorage) GETLEN() int {
	return len(product)
}
