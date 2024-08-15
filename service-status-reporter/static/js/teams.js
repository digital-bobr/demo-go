function deleteTeam(id, name) {
    if (confirm(`Are you sure you want to delete team ${name} with ID ${id}?`)) {
        console.log(`Delete team with ID: ${id}`);
        fetch(`/team/delete/${id}`, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json',
            }
        })
            .then(response => {
                if (response.ok) {
                    window.location.reload();
                } else {
                    return response.json().then(data => {
                        var errorMessage = data.error + "\n" + data.details;
                        throw new Error(errorMessage || 'Unknown error occurred');
                    });
                }
            })
            .catch(error => {
                alert(`Error: ${error.message}`)
                console.error('Error:', error);
            });
    } else {
        console.log('Deletion canceled.');
    }
}

function addNewService() {
    console.log(`Redirect to add new service page`);
    window.location.href = `/service/add`;
}

function addNewContact() {
    console.log(`Redirect to add new contact page`);
    window.location.href = `/contact/add`;
}

function addNewTeam() {
    console.log(`Redirect to add new team page`);
    window.location.href = `/team/add`;
}
