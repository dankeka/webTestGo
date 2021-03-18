const email = document.getElementById("email");
const name = document.getElementById("name");
const password1 = document.getElementById("password1");
const password2 = document.getElementById("password2");


function checkEmail() {
  let reEmail = /.+\@.+\..+/g;

  return reEmail.test(email.value);
}

function checkPassword() {
  return password1.value == password2.value;
}

document.getElementById("formRegister").onsubmit = async (e) => {
  e.preventDefault();
  
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

  return false;
};