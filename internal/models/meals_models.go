package models

import "strings"

var Vegetables = map[string]int{
	"Aceitunas negras": 349, "Aceitunas verdes": 132, "Acelgas": 33, "Ajos": 169, "Alcachofas": 64, "Apio": 20, "Berenjena": 29, "Berros": 21, "Brócoli": 31, "Calabacín": 31, "Calabaza": 24, "Cebolla": 47, "Cebolla tierna": 39, "Champiñón y otras setas": 28, "Col": 28, "Col de Bruselas": 54, "Coliflor": 30, "Endibia": 22, "Escarola": 37, "Espárragos": 26, "Espárragos en lata": 24, "Espinaca": 32, "Espinacas congeladas": 25, "Habas tiernas": 64, "Hinojo": 16, "Lechuga": 18, "Nabos": 29, "Pepino": 12, "Perejil": 55, "Pimiento": 22, "Porotos verdes": 21, "Puerros": 42, "Rábanos": 20, "Remolacha": 40, "Repollo": 19, "Rúcula": 37, "Brotes de Soja": 50, "Tomate triturado en conserva": 39, "Tomates": 22, "Trufa": 92, "Zanahoria": 42, "Zumo de tomate": 21,
}
var Fruits = map[string]int{
	"Arándanos": 41, "Caqui": 64, "Cereza": 47, "Chirimoya": 78, "Ciruela": 44, "Ciruela seca": 290, "Coco": 646, "Dátil": 279, "Dátil seco": 306, "Frambuesa": 40, "Fresas": 36, "Granada": 65, "Grosella": 37, "Higos": 80, "Higos secos": 275, "Kiwi": 51, "Limón": 39, "Mandarina": 40, "Mango": 57, "Manzana": 52, "Melón": 31, "Mora": 37, "Naranja": 44, "Nectarina": 64, "Nísperos": 97, "Papaya": 45, "Pera": 61, "Piña": 51, "Piña en almíbar": 84, "Plátano": 90, "Pomelo": 30, "Sandía": 30, "Uva": 81, "Uva pasa": 324, "Zumo de fruta": 45, "Zumo de Naranja": 42,
}
var Dairy = map[string]int{
	"Cuajada": 92, "Flan de huevo": 126, "Flan de vainilla": 102, " Helados lácteos": 167, "Leche condensada c/azúcar": 350, "Leche condensada s/azúcar": 160, "Leche de cabra": 72, "Leche de oveja": 96, "Leche descremada": 36, "Leche en polvo descremada": 373, "Leche en polvo entera": 500, "Leche entera": 68, "Leche semi descremada": 49, "Mousse": 177, "Nata o crema de leche": 298, "Queso blanco desnatado": 70, "Queso Brie": 263, "Queso cammembert": 312, "Queso cheddar": 381, "Queso crema": 245, "Queso de bola": 349, "Queso de Burgos": 174, "Queso de oveja": 245, "Queso edam": 306, "Queso emmental": 415, "Queso fundido untable": 285, "Queso gruyere": 391, "Queso manchego": 376, "Queso mozzarella": 245, "Queso parmesano": 393, "Queso ricota": 400, "Queso roquefort": 405, "Requesón": 96, "Yogur desnatado": 45, "Yogur desnatado con frutas": 82, "Yogur enriquecido con nata": 65, "Yogur natural": 62, "Yogur natural con fruta": 100,
}
var Meat = map[string]int{
	"Bacon (Panceta ahumada)": 665, "Butifarra cocida": 390, "Butifarra, salchicha fresca": 326, "Cabrito": 127, "Cerdo, chuleta": 330, "Cerdo, hígado": 153, "Cerdo, lomo": 208, "Chicharrón": 601, "Chorizo": 468, "Ciervo": 120, "Codorniz y perdiz": 114, "Conejo, liebre": 162, "Cordero lechón": 105, "Cordero pierna": 98, "Cordero, costillas": 215, "Cordero, hígado": 132, "Faisán": 144, "Foie-Gras": 518, "Gallina": 369, "Hamburguesa": 230, "Jabalí": 107, "Jamón": 380, "Jamón cocido": 126, "Jamón crudo": 296, "Jamón York": 289, "Lengua de vaca": 191, "Lomo embuchado": 380, "Mortadela": 265, "Pato": 200, "Pavo, Muslo": 186, "Pavo, Pechuga": 134, "Perdiz": 120, "Pies de cerdo": 290, "Pollo, Hígado": 129, "Pollo, Muslo": 186, "Pollo": 134, "Salami": 325, "Salchicha Frankfurt": 315, "Salchichón": 294, "Ternera": 181, "Ternera, chuleta": 168, "Ternera, hígado": 140, "Ternera, lengua": 207, "Ternera, riñón": 86, "Ternera, sesos": 125, "Ternera, solomillo": 290, "Tira de asado": 401, "Tripas": 100, "Vacuno, Hígado": 129,
}
var Fish = map[string]int{
	"Almejas": 50, "Anchoas": 175, "Anguilas": 200, "Atún en lata con aceite vegetal": 280, " Atún en lata con agua": 127, "Atún fresco": 225, "Bacalao fresco": 74, "Bacalao seco": 322, "Besugo": 118, "Caballa": 153, "Calamar": 82, "Cangrejo": 85, "Caviar": 233, "Congrio": 112, "Dorada": 80, "Gallo": 73, "Gambas": 96, "Langosta": 67, "Langostino": 96, "Lenguado": 73, "Lubina": 118, "Lucio": 81, "Mejillón": 74, "Merluza": 86, "Mero": 118, "Ostras": 80, "Pejerrey": 87, "Pez espada": 109, "Pulpo": 57, "Rodaballo": 81, "Salmón": 172, "Salmón ahumado": 154, "Salmonete": 97, "Sardina en lata con aceite vegetal": 192, "Sardinas": 151, "Trucha": 94,
}
var Cereals = map[string]int{
	"Arroz blanco": 354, "Arroz integral": 350, "Avena": 367, "Cebada": 373, "Centeno": 350, "Cereales con chocolate": 358, "Cereales desayuno, con miel": 386, "Copos de maíz": 350, "Harina de maíz": 349, "Harina de trigo integral": 340, "Harina de trigo refinada": 353, "Pan de centeno": 241, "Pan de trigo blanco": 255, "Pan de trigo integral": 239, "Pan de trigo molde blanco": 233, "Pan de trigo molde integral": 216, "Pasta al huevo": 368, "Pasta de sémola": 361, "Patatas cocidas": 86, "Patatas fritas": 312, "Polenta": 358, "Sémola de trigo": 368, "Yuca": 338,
}
var Legumes = map[string]int{
	"Garbanzos": 361, "Judías": 343, "Lentejas": 336,
}
var Eggs = map[string]int{
	"Clara": 48, "Huevo duro": 147, "Huevo entero": 162, "Yema": 368, "Huevo frito": 416,
}
var Sauce = map[string]int{
	"Bechamel": 115, "Caldos concentrados": 259, "Ketchup": 98, "Mayonesa": 718, "Mayonesa light": 374, "Mostaza": 15, "Salsa de soja": 61, "Salsa de tomate en conserva": 86, "Sofrito": 116, "Vinagres": 8,
}
var Ingredients = map[string]map[string]int{
	"Verduras":          Vegetables,
	"Frutas":            Fruits,
	"Lácteos":           Dairy,
	"Carnes":            Meat,
	"Pescados":          Fish,
	"Pastas y Cereales": Cereals,
	"Legumbres":         Legumes,
	"Huevos":            Eggs,
	"Salsas":            Sauce,
}

type MealDB struct {
	Id          string `db:"id" json:"id,omitempty"`
	UserId      string `db:"user_id" json:"userId"`
	Name        string `db:"name" json:"name" validate:"required"`
	Description string `db:"description" json:"description,omitempty"`
	Image       string `db:"image" json:"image,omitempty"`
	Type        string `db:"type" json:"type" validate:"required,oneof=semanal ocasional normal"`
	Ingredients string `db:"ingredients" json:"ingredients" validate:"required"`
	Kcal        int    `db:"kcal" json:"kcal"`
	Seasons     string `db:"seasons" json:"seasons"`
}

type Meal struct {
	Id          string   `json:"id"`
	UserId      string   `json:"user_id"`
	Name        string   `json:"name" validate:"required"`
	Description string   `json:"description"`
	Image       string   `json:"image"`
	Type        string   `json:"type" validate:"required,oneof=semanal ocasional normal"`
	Ingredients []string `json:"ingredients"`
	Kcal        int      `json:"kcal"`
	Seasons     []string `json:"seasons" validate:"required,dive,oneof=primavera verano otoño invierno general"`
}

type MealsFilters struct {
	Name    *string  `query:"name"`
	Type    *string  `query:"type"`
	Healthy *bool    `query:"healthy"`
	Season  []string `query:"[]season"`
}

func MealToAPI(meal *MealDB) *Meal {
	return &Meal{
		Id:          meal.Id,
		UserId:      meal.UserId,
		Name:        meal.Name,
		Description: meal.Description,
		Image:       meal.Image,
		Type:        meal.Type,
		Ingredients: strings.Split(meal.Ingredients, ","),
		Kcal:        meal.Kcal,
		Seasons:     strings.Split(meal.Seasons, ","),
	}
}

func MealFromAPI(meal *Meal) *MealDB {
	return &MealDB{
		Id:          meal.Id,
		UserId:      meal.UserId,
		Name:        meal.Name,
		Description: meal.Description,
		Image:       meal.Image,
		Type:        meal.Type,
		Ingredients: strings.Join(meal.Ingredients, ","),
		Kcal:        meal.Kcal,
		Seasons:     strings.Join(meal.Seasons, ","),
	}
}

type ExternalMeals struct {
	Hits []Hits `json:"hits"`
}

type Hits struct {
	Recipe Recipe `json:"recipe"`
}

type Recipe struct {
	Name        string  `json:"label"`
	Imagen      string  `json:"image"`
	Description string  `json:"url"`
	Kcal        float64 `json:"calories"`
}

type ExternalMealFilter struct {
	Q string `query:"q"`
}

func FromExternalToInternal(meal Recipe) Meal {
	return Meal{
		Name:        meal.Name,
		Image:       meal.Imagen,
		Description: meal.Description,
		Kcal:        int(meal.Kcal / 10),
		Seasons:     []string{"general"},
		Type:        "normal",
	}
}
