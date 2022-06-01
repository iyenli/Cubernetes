from serverless import Request, Response, Invoke
from http import HTTPStatus

def action(req: Request) -> Invoke or Response:
    a = int(req.param('a'))
    b = int(req.param('b'))
    s = req.param("sign")

    if s == '+':
        return Invoke("addition", params={'a': a, 'b': b})
    elif s == '*':
        return Invoke("multiplication", params={'a': a, 'b': b})
    elif s == '>':
        return Response(str(a > b) + "\n")
    elif s == '<':
        return Response(str(a < b) + "\n")
    elif s == "==":
        return Response(str(a == b) + "\n")
    else:
        return Response("unknown sign: {}\n".format(s), http_status=HTTPStatus.BAD_REQUEST)

