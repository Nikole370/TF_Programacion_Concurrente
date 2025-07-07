import { Component } from '@angular/core';
import { HttpClient } from '@angular/common/http';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  formData = {
    ubigeo: '',
    dominio: '',
    estrato: '',
    gastoG: '',
    gastoH: '',
    usoTIC: ''
  };

  resultado: any = null;
  yaEntrenado = false;

  constructor(private http: HttpClient) {}

  enviarDatos() {
    if (!this.yaEntrenado) {
      this.http.post('http://localhost:8080/train', {}).subscribe(() => {
        this.yaEntrenado = true;

        // ⏳ Esperar 7 segundos antes de predecir
        setTimeout(() => {
          this.realizarPrediccion();
        }, 7000);
      });
    } else {
      this.realizarPrediccion();
    }
  }


  realizarPrediccion() {
    this.http.post('http://localhost:8080/predict', {
      estrato: parseFloat(this.formData.estrato),
      dominio: parseFloat(this.formData.dominio)
    }).subscribe(response => {
      this.resultado = response;
    });
  }
}
