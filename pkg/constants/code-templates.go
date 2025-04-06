package constants

const PYTHON_CODE_TEMPLATE = `
import json

%s

results = []
testcases = %s

for tc in testcases:
	args = [arg.split("=")[1].strip() for arg in tc["input"].split(";")]

	inp = "("
	for arg in args:
		inp += arg + ", "
	inp += ")"

	func_name = %s

	output = eval(f"{func_name.__name__}{inp}")

	res = {
		"input": tc["input"],
		"output": output,
		"expected": tc["output"],
	}

	if f"{output}" != tc["output"]:
		res["result"] = False
	else:
		res["result"] = True
	
	results.append(res)

print(json.dumps(results))
`

const JAVASCRIPT_CODE_TEMPLATE = `
%s

const results = [];
let testcases = %s

for (const tc of testcases) {
    const args = tc.input.split(";").map(arg => arg.split("=")[1])

	let input = "("
	for(const arg of args) {
		input += arg + ", "
	}
	input += ")"

	const func_name = '%s'

	%s

	const res = {
		"input": tc.input,
		"output": output.toString(),
		"expected": tc.output,
	}
	
	if (output.toString() !== tc.output) {
		res["result"] = false
	} else {
		res["result"] = true
	}

    results.push(res);
}

console.log(JSON.stringify(results));
`

var GOLANG_CODE_TEMPLATE = map[string]string{
	"python":     PYTHON_CODE_TEMPLATE,
	"javascript": JAVASCRIPT_CODE_TEMPLATE,
}