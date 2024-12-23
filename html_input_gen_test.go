package structvalidator

import (
	"testing"
)

func TestGenerateHTML(t *testing.T) {
	s := &Test1{}
	fieldsHTMLInputs := GenerateHTML(s, &HTMLOptions{
		RestrictFields: map[string]bool{
			"FirstName": true,
			"Age":       true,
			"PostCode":  true,
			"Email":     true,
			"Country":   true,
			"County":    true,
		},
		IDPrefix:   "id_",
		NamePrefix: "name_",
	})

	if fieldsHTMLInputs["FirstName"] != `<input type="text" name="name_FirstName" id="id_FirstName" required minlength="5" maxlength="25"/>` {
		t.Fatal("GenerateHTML failed to output HTML for 'FirstName' field")
	}
	if fieldsHTMLInputs["Age"] != `<input type="number" name="name_Age" id="id_Age" required min="18" max="150"/>` {
		t.Fatal("GenerateHTML failed to output HTML for 'Age' field")
	}
	if fieldsHTMLInputs["PostCode"] != `<input type="text" name="name_PostCode" id="id_PostCode" required pattern="^[0-9][0-9]-[0-9][0-9][0-9]$"/>` {
		t.Fatal("GenerateHTML failed to output HTML for 'PostCode' field")
	}
	if fieldsHTMLInputs["Email"] != `<input type="email" name="name_Email" id="id_Email" required/>` {
		t.Fatal("GenerateHTML failed to output HTML for 'Email' field")
	}
	if fieldsHTMLInputs["Country"] != `<input type="text" name="name_Country" id="id_Country" pattern="^[A-Z][A-Z]$"/>` {
		t.Fatal("GenerateHTML failed to output HTML for 'Country' field")
	}
	if fieldsHTMLInputs["County"] != `<textarea name="name_County" id="id_County" maxlength="40"></textarea>` {
		t.Fatal("GenerateHTML failed to output HTML for 'County' field")
	}
}

func TestGenerateHTMLWithValues(t *testing.T) {
	s := &Test1{}
	fieldsHTMLInputs := GenerateHTML(s, &HTMLOptions{
		RestrictFields: map[string]bool{
			"FirstName": true,
			"Age":       true,
			"PostCode":  true,
			"Email":     true,
			"Country":   true,
			"County":    true,
		},
		IDPrefix:   "id_",
		NamePrefix: "name_",
		OverwriteValues: map[string]string{
			"FirstName": `Joe "Joe"`,
			"Age":       "40",
			"Email":     "email@example.com",
			"Country":   "XX",
		},
	})

	if fieldsHTMLInputs["FirstName"] != `<input type="text" name="name_FirstName" id="id_FirstName" required minlength="5" maxlength="25" value="Joe &#34;Joe&#34;"/>` {
		t.Fatal("GenerateHTML failed to output HTML for 'FirstName' field")
	}
	if fieldsHTMLInputs["Age"] != `<input type="number" name="name_Age" id="id_Age" required min="18" max="150" value="40"/>` {
		t.Fatal("GenerateHTML failed to output HTML for 'Age' field")
	}
	if fieldsHTMLInputs["PostCode"] != `<input type="text" name="name_PostCode" id="id_PostCode" required pattern="^[0-9][0-9]-[0-9][0-9][0-9]$"/>` {
		t.Fatal("GenerateHTML failed to output HTML for 'PostCode' field")
	}
	if fieldsHTMLInputs["Email"] != `<input type="email" name="name_Email" id="id_Email" required value="email@example.com"/>` {
		t.Fatal("GenerateHTML failed to output HTML for 'Email' field")
	}
	if fieldsHTMLInputs["Country"] != `<input type="text" name="name_Country" id="id_Country" pattern="^[A-Z][A-Z]$" value="XX"/>` {
		t.Fatal("GenerateHTML failed to output HTML for 'Country' field")
	}
	if fieldsHTMLInputs["County"] != `<textarea name="name_County" id="id_County" maxlength="40"></textarea>` {
		t.Fatal("GenerateHTML failed to output HTML for 'County' field")
	}
}

func TestGenerateHTMLWithFieldValues(t *testing.T) {
	s := &Test1{
		FirstName: "Joe",
		Age:       60,
		Email:     "joe@example.com",
	}
	fieldsHTMLInputs := GenerateHTML(s, &HTMLOptions{
		RestrictFields: map[string]bool{
			"FirstName": true,
			"Age":       true,
			"Email":     true,
		},
		IDPrefix:   "id_",
		NamePrefix: "name_",
		OverwriteValues: map[string]string{
			"FirstName": `Joe "Joe"`,
		},
		FieldValues: true,
	})

	if fieldsHTMLInputs["FirstName"] != `<input type="text" name="name_FirstName" id="id_FirstName" required minlength="5" maxlength="25" value="Joe &#34;Joe&#34;"/>` {
		t.Fatal("GenerateHTML failed to output HTML for 'FirstName' field")
	}
	if fieldsHTMLInputs["Age"] != `<input type="number" name="name_Age" id="id_Age" required min="18" max="150" value="60"/>` {
		t.Fatal("GenerateHTML failed to output HTML for 'Age' field")
	}
	if fieldsHTMLInputs["Email"] != `<input type="email" name="name_Email" id="id_Email" required value="joe@example.com"/>` {
		t.Fatal("GenerateHTML failed to output HTML for 'Email' field")
	}
}
