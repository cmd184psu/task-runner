async function fetchStatus() {
    try {
        const response = await fetch('/status');
        const data = await response.json();

        document.getElementById('progress').value = data.progress;

        updateTable('queuedTasks', data.queued);
        updateTable('runningTasks', data.running, true);
        updateTable('completedTasks', data.completed, false, true);

    } catch (error) {
        console.error('Error fetching status:', error);
    }
}

function updateTable(tableId, tasks, isRunning = false, isCompleted = false) {
    const tableBody = document.getElementById(tableId);
    tableBody.innerHTML = '';

    tasks.forEach(task => {
        const row = document.createElement('tr');
        row.innerHTML = `
            <td>${task.id}</td>
            <td>${task.command}</td>
            <td>${(isRunning || isCompleted) ? formatDuration(task.duration) : ''}</td>
            ${isCompleted ? `<td><a href="javascript:void(0);" onclick="showOutput('${task.output}')">View Output</a></td>` : ''}
        `;
        tableBody.appendChild(row);
    });
}

function formatDuration(duration) {
    const seconds = Math.floor(duration / 1000000000);
    const minutes = Math.floor(seconds / 60);
    return `${minutes}m ${seconds % 60}s`;
}

function showOutput(output) {
    alert(output);
}

async function addTask(event) {
    event.preventDefault();
    const commandInput = document.getElementById('command');
    const command = commandInput.value;
    commandInput.value = '';

    const response = await fetch('/add', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ command }),
    });

    if (response.ok) {
        fetchStatus();
    } else {
        alert('Failed to add task');
    }
}

function toggleTable(tableId) {
    const table = document.getElementById(tableId);
    if (table.classList.contains('hidden')) {
        table.classList.remove('hidden');
    } else {
        table.classList.add('hidden');
    }
}



setInterval(fetchStatus, 5000);
fetchStatus();
