function addNewTeam() {
    var name = prompt("Enter the new team name:");
    if (name) {
        if (name) {
            var formData = new FormData();
            formData.append('name', name);

            fetch('/team/create', {
                method: 'POST',
                body: formData
            })
                .then(response => response.json())
                .then(data => {
                    if (data.success) {
                        var selectElement = document.getElementById('team_dropdown_id');
                        if (selectElement) {
                            var option = document.createElement("option");
                            option.value = data.id; // Set the value to the ID returned by the server
                            option.text = name;  // Set the text to display
                            selectElement.add(option);
                            selectElement.value = data.id; // Preselect the new entry
                        }
                    } else {
                        alert('Failed to add new team');
                    }
                })
                .catch(error => console.error('Error:', error));
        }
    }
}

function addNewService() {
    console.log(`Redirect to add new service page`);
    window.location.href = `/service/add`;
}

function createNewContact() {
    console.log(`Redirect to add new contact page`);
    window.location.href = `/contact/add`;
}

function createNewTeam() {
    console.log(`Redirect to add new team page`);
    window.location.href = `/team/add`;
}
