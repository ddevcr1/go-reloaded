package modifier

import (
	"testing"
)

func TestProcessText(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		// Конвертация чисел (bin, hex)
		{
			name:     "001",
			input:    "It has been 10 (bin) years",
			expected: "It has been 2 years",
		},
		{
			name:     "002",
			input:    "I have 2A (hex) apples.",
			expected: "I have 42 apples.",
		},
		{
			name:     "003",
			input:    "1E (hex) files were added",
			expected: "30 files were added",
		},
		{
			name:     "004",
			input:    "7D (hex) files were added",
			expected: "125 files were added",
		},
		{
			name:     "005",
			input:    "Simply add 42 (hex) and 10 (bin) and you will see the result is 68.",
			expected: "Simply add 66 and 2 and you will see the result is 68.",
		},
		{
			name:     "006",
			input:    "I have to pack 101 (bin) outfits. Packed 1a (hex) just to be sure",
			expected: "I have to pack 5 outfits. Packed 26 just to be sure",
		},
		{
			name:     "007", // каскадные
			input:    "a (hex) (bin)",
			expected: "2",
		},
		{
			name:     "008", // tricky-каскад
			input:    "a(hex).(bin)",
			expected: "10.",
		},
		// Плохие комбинации с (hex) и (bin) как noise
		{
			name:     "012",
			input:    "abg (hex)",
			expected: "abg",
		},
		{
			name:     "013",
			input:    "z (hex)",
			expected: "z",
		},
		{
			name:     "014",
			input:    "12 (bin)",
			expected: "12",
		},
		{
			name:     "015",
			input:    "11080 (bin)",
			expected: "11080",
		},
		{
			name:     "016",
			input:    "AG (hex)",
			expected: "AG",
		},
		{
			name:     "017",
			input:    "XYZ (up)(hex)",
			expected: "XYZ",
		},
		// Текстовые теги (up, low, cap)
		{
			name:     "018",
			input:    "Ready, set, go (up) !",
			expected: "Ready, set, GO!",
		},
		{
			name:     "019",
			input:    "I should stop SHOUTING (low)",
			expected: "I should stop shouting",
		},
		{
			name:     "020",
			input:    "Welcome to the Brooklyn bridge (cap)",
			expected: "Welcome to the Brooklyn Bridge",
		},
		{
			name:     "021",
			input:    "This is so exciting (up, 2)",
			expected: "This is SO EXCITING",
		},
		{
			name:     "022",
			input:    "We are GOING THERE (low, 2) now.",
			expected: "We are going there now.",
		},
		{
			name:     "023",
			input:    "the great gatsby (cap, 3)",
			expected: "The Great Gatsby",
		},
		{
			name:     "024",
			input:    "ALL (LOW) (low) words are low here",
			expected: "all words are low here",
		},
		{
			name:     "025", // вложенные
			input:    "sOme WORDS (LOW, 2 (low, 2))",
			expected: "some words",
		},
		// {
		// 	name:     "026",
		// 	input:    "A an an banana a (up, 3)",
		// 	expected: "A an A BANANA A",
		// },

		// Многоступенчатые теги с вложениями
		{
			name:     "027",
			input:    "he (up, 1) said: 'i (cap) am here (up)'",
			expected: "HE said: 'I am HERE'",
		},
		{
			name:     "028",
			input:    "this is (low, 2) (cap, 2)",
			expected: "This Is",
		},
		{
			name:     "029",
			input:    "Abc (L(low)O(low)W(low))",
			expected: "Abc (l o w)",
		},
		{
			name:     "030",
			input:    "hello (low, 1)(cap, 1)(up, 1)",
			expected: "HELLO",
		},

		// Артикли a/an (включая edge cases)
		{
			name:     "031",
			input:    "There it was. A amazing rock!",
			expected: "There it was. An amazing rock!",
		},
		{
			name:     "032",
			input:    "It was a hour ago.",
			expected: "It was an hour ago.",
		},
		{
			name:     "033",
			input:    "There is no greater agony than bearing a untold story inside you.",
			expected: "There is no greater agony than bearing an untold story inside you.",
		},
		{
			name:     "034",
			input:    "a apple",
			expected: "an apple",
		},
		{
			name:     "035", // сохранение регистра
			input:    "A Apple",
			expected: "An Apple",
		},
		{
			// name:     "036", // сохранение регистра
			// input:    "A APPLE",
			// expected: "AN APPLE",
		},
		{
			name:     "037",
			input:    "an cat",
			expected: "a cat",
		},
		{
			name:     "038",
			input:    "An Cat",
			expected: "A Cat",
		},
		{
			// name:     "039",
			// input:    "a hat",
			// expected: "a hat",
		},
		{
			name:     "040",
			input:    "a onion",
			expected: "an onion",
		},
		{
			// name:     "041", // перед союзами
			// input:    "an, a and the",
			// expected: "an, a and the",
		},
		{
			name:     "042",
			input:    "I am a optimist, but a optimist",
			expected: "I am an optimist, but an optimist",
		},
		{
			name:     "043",
			input:    "harold wilson (cap, 2) : ' I am a optimist ,but a optimist who carries a raincoat . '",
			expected: "Harold Wilson: 'I am an optimist, but an optimist who carries a raincoat.'",
		},
		// Tricky case
		{
			name:     "044",
			input:    "an résumé",
			expected: "a résumé",
		},
		{
			name:     "045",
			input:    "a naïve",
			expected: "a naïve",
		},

		// Кавычки (внутренние, вложенные, cleanup)
		{
			// 	name:     "046",
			// 	input:    "'This is 'nested' quote'",
			// 	expected: "'This is nested quote'",
			// },
			// {
			// 	name:     "047",
			// 	input:    "\"This is \"nested\" quote\"",
			// 	expected: "\"This is nested quote\"",
		},
		{
			name:     "048",
			input:    "\" I can't believe \"",
			expected: "\"I can't believe\"",
		},
		// {
		// 	name:     "049",
		// 	input:    "\"Make this (up) uppercase.\"",
		// 	expected: "\"Make THIS uppercase.\"",
		// },
		// {
		// 	// name:     "050",
		// 	// input:    "\"Hello,   ,                    world!   !\"",
		// 	// expected: "\"Hello, world!!\"",
		// },
		{
			name:     "051",
			input:    "\" hello   \"     \"    there\"",
			expected: "\"hello\" \"there\"",
		},
		// {
		// 	name:     "053",
		// 	input:    "\"I am a optimist, \" he said.",
		// 	expected: "\"I am an optimist,\" he said.",
		// },

		// Fallback quotes усложнённый inWordApostrophe
		// {
		// 	name:     "054",
		// 	input:    "'kek'kek'kek'keke'kek'kek'kek'rock'n'roll'",
		// 	expected: "'kek' kek 'kek' keke 'kek' kek 'kek' rock'n'roll'",
		// },

		// Апострофы
		// {
		// 	name:     "055",
		// 	input:    "can' t",
		// 	expected: "can't",
		// },
		// {
		// 	name:     "056",
		// 	input:    "I 'm aware, but don 't worry",
		// 	expected: "I'm aware, but don't worry",
		// },
		// {
		// 	name:     "057", // Possessive Apostrophe
		// 	input:    "my systers ' room",
		// 	expected: "my systers' room",
		// },

		// // Ошибки OCR с лишними символами
		// {
		// 	name:     "058",
		// 	input:    "hello ... ...",
		// 	expected: "hello...",
		// },
		// {
		// 	name:     "059",
		// 	input:    "a amazing , , !",
		// 	expected: "an amazing!",
		// },

		// // Мультиязычный input, не как цель, а как noise
		// {
		// 	name:     "060",
		// 	input:    "a привет",
		// 	expected: "a привет",
		// },
		// {
		// 	name:     "061",
		// 	input:    "an 你好",
		// 	expected: "an 你好",
		// },
		// {
		// 	name:     "062",
		// 	input:    "an 東京",
		// 	expected: "an 東京",
		// },
		// {
		// 	name:     "063",
		// 	input:    "A 東京",
		// 	expected: "A 東京",
		// },

		// Stress-тест
		{
			name:     "064",
			input:    "a (cap, 2) (low, 2) b (up, 2)",
			expected: "A B",
		},
		// Редкие сокращения и формы
		{
			name:     "065",
			input:    "ma'am said hello",
			expected: "ma'am said hello",
		},
		// {
		// 	name:     "066",
		// 	input:    "fo'c'sle is a word",
		// 	expected: "fo'c'sle is a word",
		// },

		// Чистый флекс
		{
			name:     "068", // невалидный tag, но без слов
			input:    "(up)",
			expected: "",
		},
		{
			name:     "069", // проверка на сработку артиклей
			input:    "(aaaa)",
			expected: "(aaaa)",
		},
		{
			name:     "070", // незакрытый tag = невалидный
			input:    "word (up,101",
			expected: "word (up, 101",
		},
		{
			name:     "071", // отрицательынй tag = невалидный
			input:    "hello (cap, -2) world",
			expected: "hello world",
		},
		{
			name:     "072", // tag, 0 = невалидный
			input:    "hello (up, 0) world",
			expected: "hello world",
		},
		{
			name:     "074", // каскад валидны и невалидных
			input:    "AB (Hex)(low)",
			expected: "171",
		},
		{
			name:     "075", // каскад валидны и невалидных
			input:    "ab (up)(Hex)(low)",
			expected: "171",
		},
		// Пунктуация и пробелы
		{
			name:     "077", // двойной !! допущен в задании, хотя верно !!!
			input:    "I was sitting over there ,and then BAMM !!",
			expected: "I was sitting over there, and then BAMM!!",
		},
		{
			name:     "078",
			input:    "I was thinking ... You were right",
			expected: "I was thinking... You were right",
		},
		{
			name:     "079",
			input:    "I am exactly how they describe me: ' awesome '",
			expected: "I am exactly how they describe me: 'awesome'",
		},
		{
			name:     "080",
			input:    "As Elton John said: ' I am the most well-known homosexual in the world '",
			expected: "As Elton John said: 'I am the most well-known homosexual in the world'",
		},
		{
			name:     "081",
			input:    "I was sitting over     there !?  and then           BAMM !  !  !",
			expected: "I was sitting over there!? and then BAMM!!!",
		},
		{
			name:     "082",
			input:    "I was thinking .  .    . You were right",
			expected: "I was thinking... You were right",
		},
		{
			name:     "083",
			input:    "Punctuation tests are ... kinda boring ,what do you think !?",
			expected: "Punctuation tests are... kinda boring, what do you think!?",
		},
		{
			name:     "084",
			input:    "Hello:world.How:are you?",
			expected: "Hello: world. How: are you?",
		},
		{
			name:     "085",
			input:    "hello,there",
			expected: "hello, there",
		},
		{
			name:     "086",
			input:    "I was sitting over    !? . there",
			expected: "I was sitting over!?. there",
		},
		{
			name:     "087",
			input:    "I was thinking .  .    .",
			expected: "I was thinking...",
		},
		{
			name:     "088",
			input:    "BAMM !  !  !",
			expected: "BAMM!!!",
		},
		{
			name:     "089",
			input:    "Punctuation tests are           .    .         .",
			expected: "Punctuation tests are...",
		},
		{
			name:     "090",
			input:    "Hello:world",
			expected: "Hello: world",
		},
		{
			name:     "091",
			input:    "Don not be sad ,because sad backwards is das . And das not good",
			expected: "Don not be sad, because sad backwards is das. And das not good",
		},
		{
			name:     "092",
			input:    "das . And",
			expected: "das. And",
		},

		// Дополнительно: кавычки рядом с пунктуацией
		{
			name:     "093",
			input:    "He said: ' Hello ! '",
			expected: "He said: 'Hello!'",
		},
		// {
		// 	name:     "094",
		// 	input:    "\"Cool .\"",
		// 	expected: "\"Cool.\"",
		// },
		{
			name:     "095",
			input:    "She said: ' Hello ! '",
			expected: "She said: 'Hello!'",
		},
		// {
		// 	name:     "096",
		// 	input:    "\"Oh no .\"",
		// 	expected: "\"Oh no.\"",
		// },
		// {
		// 	name:     "097",
		// 	input:    "\" I can't believe ! \"",
		// 	expected: "\"I can't believe!\"",
		// },
		// Дополнительные кавычки и артикли (edge robustness)
		// {
		// 	name:     "098",
		// 	input:    "a 'honest' man",
		// 	expected: "an 'honest' man",
		// },
		// {
		// 	name:     "099",
		// 	input:    "a \"honest\" man",
		// 	expected: "an \"honest\" man",
		// },
		{
			name:     "100",
			input:    "'a apple'",
			expected: "'an apple'",
		},
		// {
		// 	name:     "101",
		// 	input:    "an 'user'",
		// 	expected: "a 'user'",
		// },
		// {
		// 	name:     "102",
		// 	input:    "A \"NBA\"",
		// 	expected: "An \"NBA\"",
		// },

		// Расширенный in-word apostrophe
		// {
		// 	name:     "103",
		// 	input:    "You 're not or you aren 't",
		// 	expected: "You're not or you aren't",
		// },
		// {
		// 	name:     "104",
		// 	input:    "I 'll do it, they 'll help",
		// 	expected: "I'll do it, they'll help",
		// },
		// {
		// 	name:     "105",
		// 	input:    "He 'd come if you asked",
		// 	expected: "He'd come if you asked",
		// },
		// {
		// 	name:     "106",
		// 	input:    "she said. '",
		// 	expected: "she said.'",
		// },
		{
			name:     "107", // обрамление кавычками
			input:    "hi ' hi' hi",
			expected: "hi 'hi' hi",
		},
		// {
		// 	name:     "108", // обрамление кавычками + апостроф
		// 	input:    "I ' m hi 'hi ' hi",
		// 	expected: "I'm hi 'hi' hi",
		// },
		{
			name:     "109", // одинокая кавычка
			input:    "hi 'hi",
			expected: "hi 'hi",
		},
		// {
		// 	name:     "110",
		// 	input:    "'   I can't believe   '", // обрамление кавычками + апостроф
		// 	expected: "'I can't believe'",
		// },
		{
			name:     "111", // обрамление кавычками и одинокие кавычки
			input:    "  ' hi '   hi'",
			expected: "'hi' hi'",
		},

		// Случаи с cleanup пробелов
		{
			name:     "112",
			input:    " hello  ",
			expected: "hello",
		},
		{
			name:     "113",
			input:    " awesome ",
			expected: "awesome",
		},

		// вложенный каскадных tag
		{
			name:     "114",
			input:    "me (up, 10(bin))",
			expected: "ME",
		},
		// Условные некорректные или частично обработанные конструкции
		{
			name:     "115",
			input:    "(hex",
			expected: "(hex",
		},
		{
			name:     "116",
			input:    "word (up,101",
			expected: "word (up, 101",
		},
		{
			name:     "117",
			input:    "(up)",
			expected: "",
		},
		{
			name:     "118",
			input:    "(aaaa)",
			expected: "(aaaa)",
		},

		// Дополнительные edge-кейсы с вложением тегов
		{
			name:     "119",
			input:    "a (cap, a (hex) (bin))",
			expected: "A",
		},
		{
			name:     "120",
			input:    "call an (people(up,20)) FOR (low, 100) real(cap,2)",
			expected: "call a (people) For Real",
		},
		{
			name:     "121",
			input:    "AB (Hex)(low)",
			expected: "171",
		},
		{
			name:     "122",
			input:    "ab (up)(Hex)(low)",
			expected: "171",
		},

		// Дополнительные необычные числа (двоичная система слишком длинная)
		{
			name:     "123",
			input:    "010101101001000 (bin)",
			expected: "11080",
		},

		// Повторы чисел, которые не должны обрабатываться
		{
			name:     "124",
			input:    "11080 (bin)",
			expected: "11080",
		},

		// Вложенные кавычки + cleanup
		{
			name:     "125",
			input:    "\" One ' two  '  three \"",
			expected: "\"One 'two' three\"",
		},

		// Удаление лишних пробелов между словами
		{
			name:     "126",
			input:    "Elton       John",
			expected: "Elton John",
		},

		// Сокращения или просторечие
		{
			name:     "127",
			input:    "wanna chose",
			expected: "wanna chose",
		},
		// Дополнительные флексы с некорректными/пограничными тегами
		{
			name:     "128",
			input:    "hello (up, 0) world",
			expected: "hello world",
		},
		{
			name:     "129",
			input:    "hello (cap, -2) world",
			expected: "hello world",
		},

		// Смешение стилей caps/low без реального изменения
		// {
		// 	name:     "130",
		// 	input:    "A A A A A",
		// 	expected: "A A A A A",
		// },
		// {
		// 	name:     "131",
		// 	input:    "a a a a a",
		// 	expected: "a a a a a",
		// },
		// {
		// 	name:     "132",
		// 	input:    "a or b",
		// 	expected: "a or b",
		// },
		{
			name:     "133",
			input:    "aN aNd b                                           (low, 3)",
			expected: "an and b",
		},

		// Псевдослова/аббревиатуры, важные для проверки логики a/an
		// {
		// 	name:     "134",
		// 	input:    "a fbi a x-ray a nba",
		// 	expected: "an fbi an x-ray an nba",
		// },

		// Финальные проверки robustness (ничего не должно меняться)
		// {
		// 	name:     "135",
		// 	input:    "\"I can't believe it's not butter.\"",
		// 	expected: "\"I can't believe it's not butter.\"",
		// },
		{
			name:     "136",
			input:    "ma'am said hello",
			expected: "ma'am said hello",
		},
		// чистый флекс
		// {
		// 	name:     "137",
		// 	input:    "fo 'c 'sle is a word, and Rock 'n ' Roll, or rock-n-roll",
		// 	expected: "fo'c'sle is a word, and Rock'N'Roll, or rock'n'roll",
		// },
		{
			name:     "138",
			input:    "ABC (low)(cap)",
			expected: "Abc",
		},
		{
			name:     "139",
			input:    "hello (low, 1)(cap, 1)(up, 1)",
			expected: "HELLO",
		},
		// Многоступенчатая вложенность, robustness по вложениям
		{
			name:     "140",
			input:    "a (hex) (bin)",
			expected: "2",
		},
		{
			name:     "141",
			input:    "ABC (low)(cap)",
			expected: "Abc",
		},
		{
			name:     "142",
			input:    "hello (low, 1)(cap, 1)(up, 1)",
			expected: "HELLO",
		},

		// Мультиязычные примеры
		// {
		// 	name:     "143",
		// 	input:    "привет (up)",
		// 	expected: "ПРИВЕТ",
		// },
		// {
		// 	name:     "144",
		// 	input:    "a пипл",
		// 	expected: "a пипл",
		// },
		{
			"145",
			"it (cap) was the best of times, it was the worst of times (up) , it was the age of wisdom, it was the age of foolishness (cap, 6) , it was the epoch of belief, it was the epoch of incredulity, it was the season of Light, it was the season of darkness, it was the spring of hope, IT WAS THE (low, 3) winter of despair.",
			"It was the best of times, it was the worst of TIMES, it was the age of wisdom, It Was The Age Of Foolishness, it was the epoch of belief, it was the epoch of incredulity, it was the season of Light, it was the season of darkness, it was the spring of hope, it was the winter of despair.",
		},
		{
			"146",
			"If I make you BREAKFAST IN BED (low, 3) just say thank you instead of: how (cap) did you get in my house (up, 2) ?",
			"If I make you breakfast in bed just say thank you instead of: How did you get in MY HOUSE?",
		},
		{
			"147",
			"one(low)two(cap)three(up)",
			"one Two THREE",
		},
		// {
		// 	"148",
		// 	"Can 't (up) STOP",
		// 	"CAN'T STOP",
		// },
		// {
		// 	"149",
		// 	".................word! ??????????? ::::::::::::: ;;;;;;; ;;;;,,,,,,,,,, ,,,,,!!!!!!!!!! !!!!!!!!!!!..",
		// 	"... word!?:;,!!!.",
		// },
		// {
		// 	"150",
		// 	"asgfdhjaksd ........... asgdhjksaldasd (up, 2)",
		// 	"ASGFDHJAKSD... ASGDHJKSALDASD",
		// },
		// {
		// 	"151",
		// 	"A A A A (low) (hex)",
		// 	"A A A 10",
		// },
		// {
		// 	"152",
		// 	"I 'tis fine",
		// 	"I'tis fine",
		// },
		// {
		// 	"153",
		// 	"is it? a apple",
		// 	"is it? An apple",
		// },
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := Process(tc.input)
			if result != tc.expected {
				t.Errorf("case number: %q\n For input %q\n expected. %q\n got       %q", tc.name, tc.input, tc.expected, result)
			}
		})
	}
}
