describe('button', () => {
  it('passes', () => {
    cy.visit("./codegen/codegen_js/e2e_tests/button/out.html")
    cy.get('p').contains("->").should("exist")
    cy.get('p').contains("-->").should("not.exist")
    cy.get('button').click()
    cy.get('p').contains("-->").should("exist")
    cy.get('p').contains("--->").should("not.exist")
    cy.get('button').click()
    cy.get('p').contains("--->").should("exist")
  })
})