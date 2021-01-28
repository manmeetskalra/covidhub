export default {
  register: user => {
    console.log(user);
    return fetch(process.env.REACT_APP_API_URI + "/register", {
      method: "post",
      body: JSON.stringify(user),
      headers: {
        "Content-Type": "application/json"
      }
    })
      .then(res => res.json())
      .then(data => data);
  },
  login: user => {
    console.log(user);
    return fetch(process.env.REACT_APP_API_URI + "/login", {
      method: "post",
      body: JSON.stringify(user),
      headers: {
        "Content-Type": "application/json"
      }
    })
      .then(res => res.json())
      .then(data => data);
  },
  subscribeEmail: info => {
    console.log(info);
    return fetch(process.env.REACT_APP_EMAIL + "/subscribe", {
      method: "post",
      body: JSON.stringify(info),
      headers: {
        "Content-Type": "application/json"
      }
    })
      .then(res => { return res })
      .catch(err => { console.log(err) });
  },
  subscribeNowEmail: info => {
    console.log(info);
    return fetch(process.env.REACT_APP_EMAIL + "/subscribe/now", {
      method: "post",
      body: JSON.stringify(info),
      headers: {
        "Content-Type": "application/json"
      }
    })
    .then(res => { return res })
    .catch(err => { console.log(err) });
  },
  subscribeText: info => {
    console.log(info);
    return fetch(process.env.REACT_APP_TEXT + "/subscribe", {
      method: "post",
      body: JSON.stringify(info),
      headers: {
        "Content-Type": "application/json"
      }
    })
    .then(res => { return res })
    .catch(err => { console.log(err) });
  },
  subscribeNowText: info => {
    console.log(info);
    return fetch(process.env.REACT_APP_TEXT + "/subscribe/now", {
      method: "post",
      body: JSON.stringify(info),
      headers: {
        "Content-Type": "application/json"
      }
    })
    .then(res => { return res })
    .catch(err => { console.log(err) });
  },
  unsubscribeEmail: info => {
    console.log(info);
    return fetch(process.env.REACT_APP_EMAIL + "/unsubscribe", {
      method: "delete",
      body: JSON.stringify(info),
      headers: {
        "Content-Type": "application/json"
      }
    })
    .then(res => { return res })
    .catch(err => { console.log(err) });
  },
  unsubscribeText: info => {
    console.log(info);
    return fetch(process.env.REACT_APP_TEXT + "/unsubscribe", {
      method: "delete",
      body: JSON.stringify(info),
      headers: {
        "Content-Type": "application/json"
      }
    })
    .then(res => { return res })
    .catch(err => { console.log(err) });
  }
};
