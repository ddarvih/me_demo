"use strict";

const ExprBase = (base, evaluate, toString, prefix) => {
    base.prototype.evaluate = evaluate;
    base.prototype.toString = toString;
    base.prototype.prefix = prefix;
    return base;
};

const Const = ExprBase(
    function (value) {
        this.value = value;
    },
    function () {
        return this.value;
    },
    function () {
        return this.value.toString();
    },
    function () {
        return this.value.toString();
    }
);

const vars = {
    "x": (...args) => args[0],
    "y": (...args) => args[1],
    "z": (...args) => args[2],
    "t": (...args) => args[3]
};

const Variable = ExprBase(
    function (name) {
        this.name = name;
    },
    function (...args) {
        return vars[this.name](...args);
    },
    function () {
        return this.name;
    },
    function () {
        return this.name;
    }
);

const Action = ExprBase(
    function (f, symb, args) {
        this.f = f;
        this.symb = symb;
        this.args = args;
    },
    function (...argsSec) {
        return this.f(...this.args.map(a => a.evaluate(...argsSec)));
    },
    function () {
        return `${this.args.map(a => a.toString()).join(" ")} ${this.symb}`;
    },
    function () {
        return `(${this.symb} ${this.args.map(a => a.prefix()).join(" ")})`;
    }
);

const ActionFabric = (f, symb) => {
    function Operation(...args) {
        Action.call(this, f, symb, args);
    }

    Operation.prototype = Object.create(Action.prototype);
    return Operation;
};

const Add = ActionFabric((a, b) => a + b, "+");
const Subtract = ActionFabric((a, b) => a - b, "-");
const Multiply = ActionFabric((a, b) => a * b, "*");
const Divide = ActionFabric((a, b) => a / b, "/");
const Negate = ActionFabric(a => -a, "negate");
const Power = ActionFabric(Math.pow, "pow");
const Log = ActionFabric((a, b) => Math.log(Math.abs(b)) / Math.log(Math.abs(a)), "log");

const Clamp = ActionFabric((a, mi, ma) => a < mi ? mi : (a > ma ? ma : a), "clamp");
const SoftClamp = ActionFabric((x, mi, ma, l) => mi + (ma - mi) / (1 + Math.exp(l * ((ma + mi) / 2 - x))),
    "softClamp");

const SumCb = ActionFabric((...args) => args.reduce((a, b) => a + b * b * b, 0), "sumCb");
const MeanCb = ActionFabric((...args) => args.reduce((a, b) => a + b * b * b, 0) / args.length, "meanCb");

const operations = {
    "+": Add,
    "-": Subtract,
    "*": Multiply,
    "/": Divide,
    "negate": Negate,
    "pow": Power,
    "log": Log,
    "clamp": Clamp,
    "softClamp": SoftClamp,
    "sumCb": SumCb,
    "meanCb": MeanCb
};

const noLimit = -1;
const amountArgs = {
    "+": 2,
    "-": 2,
    "*": 2,
    "/": 2,
    "negate": 1,
    "pow": 2,
    "log": 2,
    "clamp": 3,
    "softClamp": 4,
    "meanCb": noLimit,
    "sumCb": noLimit
}

const parseOperation = (cur, stack) => {
    let amArgs = amountArgs[cur];
    const args = stack.slice(-amArgs);
    const operation = new operations[cur](...args);
    return [...stack.slice(0, -amArgs), operation];
};

const parseOperand = (cur, stack) => {
    if (cur in vars) {
        return [...stack, new Variable(cur)];
    }
    return [...stack, new Const(+cur)];
};

const parse = (expr) => {
    return expr.trim().split(/\s+/).filter(piece => piece !== "")
        .reduce((stack, cur) => (operations[cur] ? parseOperation(cur, stack) : parseOperand(cur, stack)), [])[0];
};

function CustomError(message) {
    this.message = message;
    this.name = "CustomError";
}
CustomError.prototype = Object.create(Error.prototype);
CustomError.prototype.constructor = CustomError;

const CustomErrorFabric = (name) => {
    function NewException(message) {
        CustomError.call(this, message);
        this.name = name;
    }
    NewException.prototype = Object.create(CustomError.prototype);
    NewException.prototype.constructor = NewException;
    return NewException;
};

const GivenExtraException = CustomErrorFabric("GivenExtraError");
const GivenLessException = CustomErrorFabric("GivenLessError");
const GivenUnknownException = CustomErrorFabric("GivenUnknownError");

const parsePrefix = (str) => {
    const tokens = str.replace(/([()])/g, " $1 ").split(/\s+/).filter(t => t !== "");
    let pos = 0;

    const parseOperationInner = (cur) => {
        let amArgs = amountArgs[cur];
        let args = []
        if (amArgs === noLimit) {
            let curArg = innerParse();
            while (curArg !== ")") {
                args.push(curArg);
                curArg = innerParse();
            }
            pos--;
            return new operations[cur](...args);
        }

        for (let i = 0; i < amArgs; i++) {
            args.push(innerParse());
        }
        return new operations[cur](...args);
    };

    const checkEnd = () => {
        if (pos >= tokens.length) {
            throw new GivenLessException("Unexpected end of input");
        }
    }

    const innerParse = () => {
        checkEnd();
        let cur = tokens[pos++];

        if (cur === "(") {
            checkEnd();
            let inBrackets = tokens[pos++];
            if (!inBrackets || !(inBrackets in operations)) {
                throw new GivenLessException("Expected operator after '(' pos: " + pos);
            }

            let resInner = parseOperationInner(inBrackets);

            if (pos >= tokens.length || tokens[pos] !== ")") {
                throw new GivenLessException("Missing ')' in pos: " + pos);
            }
            pos++;
            return resInner;
        }

        if (cur === ")") {
            return cur;
        }

        if (cur in operations) {
            throw new GivenLessException("Expected '(' before operator, pos: " + pos);
        }

        if (cur in vars) {
            return new Variable(cur);
        }

        if (/^-?\d+(\.\d+)?$/.test(cur)) {
            return new Const(+cur);
        }

        throw new GivenUnknownException("Unknown element in pos: " + pos);
    }

    let result = innerParse();
    if (pos !== tokens.length) {
        throw new GivenExtraException("Unexpected tokens after pos: " + pos);
    }
    return result;
}