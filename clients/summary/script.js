const HOST_NAME = window.config?.remoteHost
if (!HOST_NAME) {
  console.error('No host name found in config')
}
console.log(`Using summarizer host name: ${HOST_NAME}`)


const btn = document.querySelector('button')

btn.addEventListener('click', async () => {
  const input = document.querySelector('input')
  const url = input.value
  const reqUrl = `${HOST_NAME}/v1/summary?url=${url}`
  let metadata = {}
  try {
    const res = await fetch(reqUrl, {
      method: 'GET',
      mode: 'cors',
      headers: {
        'Content-Type': 'application/json'
      }
    })
    metadata = await res.json()
  } catch (e) {
    renderError(e)
    return
  }
  console.log(metadata)
  renderSummary(metadata)
})

function renderSummary (data) {
  const results = document.getElementById('results')
  results.innerHTML = ''
  if (data.title) {
    const title = document.createElement('div')
    title.innerHTML = data.title
    results.appendChild(title)
  }
  if (data.description) {
    const desc = document.createElement('div')
    desc.innerHTML = data.description
    results.appendChild(desc)
  }
  if (data.images) {
    for (const imgData of data.images) {
      const img = document.createElement('img')
      img.setAttribute('src', imgData.url)
      img.setAttribute('alt', imgData.alt)
      results.appendChild(img)
    }
  }
  if (data.videos) {
    for (const vidData of data.videos) {
      let vid
      if (vidData.type == 'text/html') vid = document.createElement('iframe')
      else if (vidData.type.startsWith('video/')) {
        vid = document.createElement('video')
      } else {
        console.log(`can't render video of type ${vidData.type}`)
        continue
      }
      vid.setAttribute('src', vidData.url)
      results.appendChild(vid)
    }
  }
}

function renderError (e) {
  const results = document.getElementById('results')
  results.innerHTML = `Error parsing metadata of given url: ${e}`
}
