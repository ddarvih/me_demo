"use strict";

const cnst = (val) => () => val;

const tau = cnst(Math.PI * 2);
const phi = cnst(((1 + Math.sqrt(5)) / 2));

const vars = {
    "x": (...args) => args[0],
    "y": (...args) => args[1],
    "z": (...args) => args[2],
    "t": (...args) => args[3]
}

const variable = (name) => (...args) => vars[name](...args);

const action = f => (...args) => (...argsSec) => f(...args.map(a => a(...argsSec)));

const add = action((a, b) => a + b);
const subtract = action((a, b) => a - b);
const multiply = action((a, b) => a * b);
const divide = action((a, b) => a / b);
const negate = action(a => -a);

const power = action(Math.pow);
const log = action((a, b) => Math.log(Math.abs(b)) / Math.log(Math.abs(a)));
const cosh = action(Math.cosh);
const sinh = action(Math.sinh);

//let myexpr = add(cnst(1),subtract(multiply(variable("x"), variable("x")), multiply(cnst(2), variable("x"))));
//for (let i = 0; i < 11; i++) {
//    console.log(myexpr(i));
//}

