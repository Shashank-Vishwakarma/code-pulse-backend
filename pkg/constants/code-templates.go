package constants

const PYTHON_CODE_TEMPLATE = `
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

print(results)
`

const JAVASCRIPT_CODE_TEMPLATE = `
%s

let results = []
let testcases = %s

for(let tc of testcases) {
	let args = tc["input"].split(";").map(arg => arg.trim())

	let output = %s + "("
	for(let arg of args):
		output += arg + ", "
	output += ")"

	const res = {
		"input": tc["input"],
		"output": output,
		"expected": tc["output"],
	}

	if output != tc["output"]:
		res["result"] = false
	else:
		res["result"] = true
	
	results.push(res)
}

console.log(results)
`

var GOLANG_CODE_TEMPLATE = map[string]string{
	"python":     PYTHON_CODE_TEMPLATE,
	"javascript": JAVASCRIPT_CODE_TEMPLATE,
}