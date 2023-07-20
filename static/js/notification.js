let toastBox = document.getElementById('toastBox')

async function clearSession(type) {
    let formData = new FormData();
    let tkn = document.querySelector('meta[name="csrf_token"]').content
    formData.append("csrf_token", tkn)
    const resp = await fetch(`/api/clear/${type}`, {
        method: "POST",
        body: formData,
    });
}

function showToast(messages, types) {
    for (let i = 0; i < messages.length; i++) {
        let toast = document.createElement('div');
        toast.classList.add('toast');
        let msg = messages[i];
        let type = types[i];
    
        if (type == "flash") {
          toast.innerHTML = `<i class="fa-solid fa-circle-check" style="color: #2bff00;"></i> ${msg}`;
        } else if (type == "error") {
          toast.innerHTML = `<i class="fa-solid fa-circle-xmark" style="color: #ff0000;"></i> ${msg}`;
        } else if (type == "warning") {
          toast.innerHTML = `<i class="fa-solid fa-triangle-exclamation" style="color: #ffc800;"></i> ${msg}`;
        }
    
        toastBox.appendChild(toast);
    
        if (type.includes('error')) {
          toast.classList.add('error');
        }
        if (type.includes('warning')) {
          toast.classList.add('warning');
        }
        clearSession(type);
        setTimeout(() => {
          toast.remove();
        }, 2000);
      }
}

if (document.querySelector('meta[name="message"]')) {
    // message = document.querySelector('meta[name="message"]').content
    // type = document.querySelector('meta[name="type"]').content
    const messages = Array.from(document.querySelectorAll('meta[name="message"]')).map(meta => meta.content);
    const types = Array.from(document.querySelectorAll('meta[name="type"]')).map(meta => meta.content);
    showToast(messages, types)
}