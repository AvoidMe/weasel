import ast
import textwrap

import pytest

from tests.common import compare_ast
from weasel.const_counter import ConstCounter
from weasel.transformers.rule_1 import Rule_1


@pytest.mark.parametrize(
    "input, expected",
    [
        ("x = 5", ""),
        ("x = 5; print(x)", "print(5)"),
        ("x = 5; y = 6; print(x + y)", "print(5 + 6)"),
        ("x = 5; x = 6; print(x)", "x = 5; x = 6; print(x)"),
        (
            textwrap.dedent(
                """
            def main():
                x = 5
                print(x)
            """
            ),
            textwrap.dedent(
                """
            def main():
                print(5)
            """
            ),
        ),
        (
            textwrap.dedent(
                """
            x = 5
            def main():
                print(x)
            """
            ),
            textwrap.dedent(
                """
            def main():
                print(5)
            """
            ),
        ),
        (
            textwrap.dedent(
                """
            x = 5
            def main():
                x = 10
                print(x)
            """
            ),
            textwrap.dedent(
                """
            def main():
                print(10)
            """
            ),
        ),
        (
            textwrap.dedent(
                """
            x = 5
            async def main():
                print(x)
            """
            ),
            textwrap.dedent(
                """
            async def main():
                print(5)
            """
            ),
        ),
        (
            textwrap.dedent(
                """
            x = 10
            def main():
                x = 5
                print(x)
            print(x)
            """
            ),
            textwrap.dedent(
                """
            def main():
                print(5)
            print(10)
            """
            ),
        ),
        (
            textwrap.dedent(
                """
            def main():
                x = 5
                print(x)
            print(x)
            """
            ),
            textwrap.dedent(
                """
            def main():
                print(5)
            print(x)
            """
            ),
        ),
        (
            textwrap.dedent(
                """
            if something:
                x = 5
            print(x)
            """
            ),
            textwrap.dedent(
                """
            if something:
                x = 5
            print(x)
            """
            ),
        ),
        (
            textwrap.dedent(
                """
            for i in range(5):
                x = 5
            print(x)
            """
            ),
            textwrap.dedent(
                """
            for i in range(5):
                x = 5
            print(x)
            """
            ),
        ),
        (
            textwrap.dedent(
                """
            while something:
                x = 5
            print(x)
            """
            ),
            textwrap.dedent(
                """
            while something:
                x = 5
            print(x)
            """
            ),
        ),
        (
            textwrap.dedent(
                """
            x = 10
            while something:
                x = 5
            print(x)
            """
            ),
            textwrap.dedent(
                """
            x = 10
            while something:
                x = 5
            print(x)
            """
            ),
        ),
        (
            textwrap.dedent(
                """
            x = 10
            for i in range(5):
                x = 5
            print(x)
            """
            ),
            textwrap.dedent(
                """
            x = 10
            for i in range(5):
                x = 5
            print(x)
            """
            ),
        ),
        (
            textwrap.dedent(
                """
            x = 10
            if something:
                x = 5
            print(x)
            """
            ),
            textwrap.dedent(
                """
            x = 10
            if something:
                x = 5
            print(x)
            """
            ),
        ),
        (
            textwrap.dedent(
                """
            def main():
                x = 5
            """
            ),
            textwrap.dedent(
                """
            def main():
                pass
            """
            ),
        ),
        (
            textwrap.dedent(
                """
            x = 5
            def main():
                global x
                x = 10
            print(x)
            """
            ),
            textwrap.dedent(
                """
            x = 5
            def main():
                global x
                x = 10
            print(x)
            """
            ),
        ),
        (
            textwrap.dedent(
                """
            x = 5
            def main():
                nonlocal x
                x = 10
            print(x)
            """
            ),
            textwrap.dedent(
                """
            x = 5
            def main():
                nonlocal x
                x = 10
            print(x)
            """
            ),
        ),
        (
            textwrap.dedent(
                """
            x = 5
            match something:
                case "a":
                    x = 10
            print(x)
            """
            ),
            textwrap.dedent(
                """
            x = 5
            match something:
                case "a":
                    x = 10
            print(x)
            """
            ),
        ),
        (
            textwrap.dedent(
                """
            x = 5
            class Something:
                def foo(self):
                    x = 10
            print(x)
            """
            ),
            textwrap.dedent(
                """
            class Something:
                def foo(self):
                    pass
            print(5)
            """
            ),
        ),
        (
            textwrap.dedent(
                """
            x = 5
            class Something:
                def foo(self):
                    global x
                    x = 10
            print(x)
            """
            ),
            textwrap.dedent(
                """
            x = 5
            class Something:
                def foo(self):
                    global x
                    x = 10
            print(x)
            """
            ),
        ),
        (
            textwrap.dedent(
                """
            class Something:
                def __init__(self):
                    self.x = 5

                def foo(self):
                    print(self.x)
            """
            ),
            textwrap.dedent(
                """
            class Something:
                def __init__(self):
                    self.x = 5

                def foo(self):
                    print(self.x)
            """
            ),
        ),
    ],
    ids=[
        "Constant without usage should be removed",
        "Single constant",
        "Two constants",
        "Shouldn't optimize, value overriden",
        "Inside function",
        "Inside function, global constant",
        "Inside function, global constant, local constant overrides",
        "Inside async function",
        "Global + local optimization",
        "Local function constant shouldn't pollute global namespace",
        "Constants inside if branches should be ignored",
        "Constants inside for loops should be ignored",
        "Constants inside while loops should be ignored",
        "Non-constants inside while loops should be ignored",
        "Non-constants inside for loops should be ignored",
        "Non-constants inside if branches should be ignored",
        "We shouldn't leave functions body empty",
        "We should ignore constants which are involved in global statement",
        "We should ignore constants which are involved in nonlocal statement",
        "Constants inside match statements should be ignored",
        "Class with same constant name",
        "Class with global statement for same constant name",
        "Class with local constant shouldn't be touched",
    ],
)
def test_constant_folding_rule(input, expected):
    tree = ast.parse(input)
    counter = ConstCounter()
    counter.visit(tree)
    result_tree = Rule_1([counter]).visit(tree)
    expected = ast.parse(expected)
    assert compare_ast(result_tree, expected)
