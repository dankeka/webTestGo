const username = document.getElementById("name");
const email = document.getElementById("email");
const password1 = document.getElementById("password1");
const password2 = document.getElementById("password2");


function checkEmail() {
  let reEmail = /.+\@.+\..+/g;

  let checkEm = reEmail.test(email.value);

  if(checkEm || !email.value) {
    return true;
  } else {
    return false;
  }
}

function checkPassword() {
  return password1.value == password2.value && password1.value != "" && password2.value != "";
}

function checkUserName() {
  return username.value != "";
}

document.getElementById("formRegister").onsubmit = async (e) => {
  if(!checkEmail()) {
    try {
      let el = document.getElementById("errEmailDiv");
      el.remove();
    } catch {
      // pass
    }

    let emailError = document.createElement("div");
    emailError.id = "errEmailDiv";
    emailError.innerHTML = `
      <span style="color: red; font-size: .9em;">Проверьте email</span>
    `;

    email.after(emailError);
  } else {
    try {
      let el = document.getElementById("errEmailDiv");
      el.remove();
    } catch {
      // pass
    }
  }

  if(!checkPassword()) {
    try {
      let el = document.getElementById("errPasswordDiv");
      el.remove();
    } catch {
      // pass
    }

    let passwordError = document.createElement("div");
    passwordError.id = "errPasswordDiv";
    passwordError.innerHTML = `
        <span style="color: red; font-size: .9em;">Пароли не совпадают</span>
    `;

    password2.after(passwordError);
  } else {
    try {
      let el = document.getElementById("errPasswordDiv");
      el.remove();
    } catch {
      // pass
    }
  }

  if(!checkUserName()) {
    try {
      let el = document.getElementById("errUserNameDiv");
      el.remove();
    } catch {
      // pass
    }

    let userNameError = document.createElement("div");
    userNameError.id = "errUserNameDiv";
    userNameError.innerHTML = `
        <span style="color: red; font-size: .9em;">Поле имени пустое!</span>
    `;

    username.after(userNameError);
  } else {
    try {
      let el = document.getElementById("errUserNameDiv");
      el.remove();
    } catch {
      // pass
    }
  }

  if(checkPassword() && checkEmail() && checkUserName()) {
    return true;
  } else {
    e.preventDefault();
    return false;
  }
};