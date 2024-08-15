function deleteService(id) {
    if (confirm(`Are you sure you want to delete service with ID ${id}?`)) {
        console.log(`Delete service with ID: ${id}`);
        fetch(`/service/delete/${id}`, {
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
            })
    } else {
        // Do nothing if user cancels
        console.log('Deletion canceled.');
    }
}

function editService(id) {
    console.log(`Edit service with ID: ${id}`);
    window.location.href = `/service/edit?id=${id}`;
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

function toggleNested(id) {
    var allNested = document.getElementsByClassName('nested');
    for (var i = 0; i < allNested.length; i++) {
        allNested[i].style.display = 'none';
    }

    var element = document.getElementById("nested-" + id);
    if (element.style.display === "none") {
        element.style.display = "table-row";
        checkHealth(id);
    } else {
        element.style.display = "none";
    }
}

function checkHealth(id, prodEnv, preprodEnv, nonprodEnv) {
    console.log(`Checking health for service ID: ${id}`);
    fetch(`/service/health/${id}`, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
        }
    })
        .then(response => {
            if (response.ok) {
                return response.json();
            } else {
                return response.json().then(data => {
                    var errorMessage = data.error + "\n" + data.details;
                    throw new Error(errorMessage || 'Unknown error occurred');
                });
            }
        })
        .then(serviceHealth => {
            updateStatus(`status-prod-nested-${id}`, serviceHealth.production);
            updateStatus(`status-preprod-nested-${id}`, serviceHealth.preproduction);
            updateStatus(`status-nonprod-nested-${id}`, serviceHealth.nonproduction);
        })
        .catch(error => {
            alert(`Error: ${error.message}`)
            console.error('Error:', error);
        })
}

function updateStatus(elementId, environment) {
    var statusElement = document.getElementById(elementId);
    if (statusElement) {
        if (environment.error) {
            statusElement.innerText = "DOWN: " + environment.error;
            statusElement.className = "status-down";
            console.error(`Error updating status for ${elementId}:`, environment.error);
        } else {
            statusElement.innerText = `Status: UP, Version: ${environment.version}, Ref: ${environment.ref}, Commit: ${environment.commit}`;
            statusElement.className = "status-up";
        }
    } else {
        console.error(`Status element not found for id: ${elementId}`);
    }
}

function showContactDetails(id, role, name, email, slack) {
    var modal = document.getElementById("contactModal");
    var contactRole = document.getElementById("contactRole");
    var contactDetails = document.getElementById("contactDetails");
    contactRole.innerText = role;
    contactDetails.innerHTML = `${name}<br><a href="mailto:${email}">${email}</a><br>${slack}`;
    modal.style.display = "block";
}

function closeModal() {
    var modal = document.getElementById("contactModal");
    modal.style.display = "none";
}

window.onclick = function (event) {
    var modal = document.getElementById("contactModal");
    if (event.target == modal) {
        modal.style.display = "none";
    }
}
