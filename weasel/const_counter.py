import ast


class NodeCounter:
    def __init__(self, node, value):
        self.node = node
        self.value = value
        self.count = 0


class ConstCounter(ast.NodeVisitor):
    def __init__(self):
        super().__init__()
        self._statistic = {}
        self.inner = 0
        self.count = [1]

    def increment_node(self, name, node, value=None, count=1):
        counter = self._statistic.get(name, NodeCounter(node, value))
        counter.count += count
        self._statistic[name] = counter

    def inner_visit(self, node):
        # we care only about global/nonlocal statements
        self.inner += 1
        for child in node.body:
            self.visit(child)
        self.inner -= 1

    def double_visit(self, node):
        """
        We should ignore everything inside for, while, if statements
        But at the same time we should include them in final count
        For cases like this:
            if something:
                x = 5
            x = 10
        """
        self.count.append(self.count[-1] + 1)
        for child in node.body:
            self.visit(child)
        self.count.pop()

    def visit_Global(self, node):
        for name in node.names:
            self.increment_node(name, node, count=2)

    def visit_Nonlocal(self, node):
        # TODO: in the future we want to trace this names to parent scope
        #       instead of complete ignore
        for name in node.names:
            self.increment_node(name, node, count=2)

    def visit_While(self, node):
        self.double_visit(node)

    def visit_For(self, node):
        self.double_visit(node)

    def visit_If(self, node):
        self.double_visit(node)

    def visit_FunctionDef(self, node):
        self.inner_visit(node)

    def visit_ClassDef(self, node):
        self.inner_visit(node)

    def visit_Assign(self, node):
        if self.inner:
            return
        # TODO: we want to optimize such cases also
        if len(node.targets) > 1 or not isinstance(node.targets[0], ast.Name):
            return
        target_name = node.targets[0].id
        if isinstance(node.value, ast.Constant):
            self.increment_node(
                target_name, node, value=node.value, count=self.count[-1]
            )
