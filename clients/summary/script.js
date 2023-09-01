HOST_NAME = 'http://localhost:4000'

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
})

function renderError (e) {
  const results = document.getElementById('results')
  results.innerHTML = `Error parsing metadata of given url: ${e}`
}
