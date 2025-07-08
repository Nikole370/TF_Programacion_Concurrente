import { Component } from '@angular/core';
import { HttpClient } from '@angular/common/http';

interface Registro {
  dominio: number;
  estrato: number;
  prediccion: number;
  resultado: string;
}

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  formData = {
    dominio: 1,
    estrato: 1
  };

  resultado: any = null;
  historial: Registro[] = [];
  yaEntrenado = false;

  constructor(private http: HttpClient) {}

  enviarDatos() {
    if (!this.yaEntrenado) {
      this.http.post('http://localhost:8080/train', {}).subscribe(() => {
        this.yaEntrenado = true;

        setTimeout(() => {
          this.realizarPrediccion();
        }, 7000);
      });
    } else {
      this.realizarPrediccion();
    }
  }

  realizarPrediccion() {
    this.http.post<any>('http://localhost:8080/predict', this.formData).subscribe(
      (res) => {
        const valor = parseFloat(res.prediccion);
        const resultadoFinal = valor >= 0.5 ? 'Usa' : 'No usa';

        this.resultado = resultadoFinal;

        this.historial.push({
          dominio: this.formData.dominio,
          estrato: this.formData.estrato,
          prediccion: valor,
          resultado: resultadoFinal
        });
      },
      (err) => {
        console.error('Error en la predicción:', err);
      }
    );
  }
}
